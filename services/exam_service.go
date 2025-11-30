package services

import (
	"HOSEROF_API/config"
	"HOSEROF_API/models"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateExam(exam models.Exam, questions []models.Question) (string, error) {
	ctx := context.Background()
	exams := config.DB.Collection("exams")
	doc := exams.NewDoc()
	exam.ExamID = doc.ID
	exam.CreatedAt = time.Now()
	exam.Released = false

	_, err := doc.Set(ctx, exam)
	if err != nil {
		return "", err
	}

	for _, q := range questions {
		if q.QID == "" {
			q.QID = doc.Collection("questions").NewDoc().ID
		}
		_, err := doc.Collection("questions").Doc(q.QID).Set(ctx, q)
		if err != nil {
			return "", err
		}
	}

	return doc.ID, nil
}

func GetExamsForClass(class string, studentID string) ([]models.Exam, error) {
	ctx := context.Background()

	q := config.DB.Collection("exams").
		Where("class", "==", class)
	iter := q.Documents(ctx)
	var out []models.Exam
	now := time.Now()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var e models.Exam
		if err := doc.DataTo(&e); err != nil {
			return nil, err
		}

		if now.Before(e.StartTime) {
			continue
		}
		if now.After(e.EndTime) {
			continue
		}

		subSnap, err := doc.Ref.Collection("submissions").Doc(studentID).Get(ctx)
		if err == nil && subSnap.Exists() {
			continue
		}

		out = append(out, e)
	}

	return out, nil
}

func GetExamQuestions(examID string, forStudent bool) ([]models.Question, error) {
	ctx := context.Background()

	qIter := config.DB.Collection("exams").
		Doc(examID).
		Collection("questions").
		OrderBy("qid", firestore.Asc).
		Documents(ctx)

	var qs []models.Question

	for {
		doc, err := qIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var q models.Question
		if err := doc.DataTo(&q); err != nil {
			return nil, err
		}

		if forStudent {
			q.CorrectAnswer = ""
		}

		qs = append(qs, q)
	}

	return qs, nil
}

func SubmitExam(examID string, studentID string, answers map[string]models.Answer) error {
	ctx := context.Background()
	examDoc := config.DB.Collection("exams").Doc(examID)

	snap, err := examDoc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errors.New("exam not found")
		}
		return err
	}
	var exam models.Exam
	if err := snap.DataTo(&exam); err != nil {
		return err
	}

	now := time.Now()
	if now.Before(exam.StartTime) {
		return errors.New("exam not started")
	}
	if now.After(exam.EndTime) {
		return errors.New("exam ended")
	}

	subDoc := examDoc.Collection("submissions").Doc(studentID)
	existsSnap, err := subDoc.Get(ctx)
	if err == nil && existsSnap.Exists() {
		return errors.New("already submitted")
	}

	qIter := examDoc.Collection("questions").Documents(ctx)
	correctMap := make(map[string]models.Question)
	for {
		doc, err := qIter.Next()
		if err != nil {
			break
		}
		var q models.Question
		if err := doc.DataTo(&q); err == nil {
			correctMap[q.QID] = q
		}
	}

	var autoScore float64 = 0
	var manualScore float64 = 0
	for qid, ans := range answers {
		q, ok := correctMap[qid]
		if !ok {
			continue
		}
		if q.Type == models.MCQ || q.Type == models.TF {
			studentAns := strings.TrimSpace(strings.ToLower(fmt.Sprint(ans.Response)))
			correctAns := strings.TrimSpace(strings.ToLower(fmt.Sprint(q.CorrectAnswer)))

			if studentAns == correctAns {
				autoScore += q.Points
				a := ans
				a.AutoScore = q.Points
				answers[qid] = a
			} else {
				a := ans
				a.AutoScore = 0
				answers[qid] = a
			}
		} else {
			a := ans
			a.AutoScore = 0
			answers[qid] = a
		}
	}

	final := autoScore + manualScore

	submission := models.Submission{
		StudentID:   studentID,
		StartedAt:   now,
		SubmittedAt: now,
		Answers:     answers,
		AutoScore:   autoScore,
		ManualScore: manualScore,
		FinalScore:  final,
		Graded:      false,
		Released:    false,
	}

	_, err = subDoc.Set(ctx, submission)
	if err != nil {
		return err
	}

	return nil
}

func GetSubmission(examID, studentID string) (models.Submission, error) {
	ctx := context.Background()
	doc := config.DB.Collection("exams").Doc(examID).Collection("submissions").Doc(studentID)
	snap, err := doc.Get(ctx)
	if err != nil {
		return models.Submission{}, err
	}
	var s models.Submission
	if err := snap.DataTo(&s); err != nil {
		return models.Submission{}, err
	}
	return s, nil
}

func GetAllSubmissions(examID string) ([]models.Submission, error) {
	ctx := context.Background()

	iter := config.DB.Collection("exams").
		Doc(examID).
		Collection("submissions").
		Documents(ctx)

	var out []models.Submission

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var s models.Submission
		if err := doc.DataTo(&s); err != nil {
			return nil, err
		}

		out = append(out, s)
	}

	return out, nil
}

func GradeWrittenAnswer(examID, studentID, qid string, score float64) error {
	ctx := context.Background()
	subRef := config.DB.Collection("exams").Doc(examID).Collection("submissions").Doc(studentID)

	snap, err := subRef.Get(ctx)
	if err != nil {
		return err
	}
	var s models.Submission
	if err := snap.DataTo(&s); err != nil {
		return err
	}

	ans, ok := s.Answers[qid]
	if !ok {
		return errors.New("question not found in submission")
	}

	ans.ManualScore = score
	s.Answers[qid] = ans

	var manualSum float64
	for _, a := range s.Answers {
		manualSum += a.ManualScore
	}
	s.ManualScore = manualSum
	s.FinalScore = s.AutoScore + s.ManualScore
	s.Graded = true

	_, err = subRef.Set(ctx, s, firestore.MergeAll)
	return err
}

func ReleaseResults(examID string) error {
	ctx := context.Background()
	examRef := config.DB.Collection("exams").Doc(examID)
	_, err := examRef.Update(ctx, []firestore.Update{
		{Path: "released", Value: true},
	})
	if err != nil {
		return err
	}

	subs := examRef.Collection("submissions").Documents(ctx)
	for {
		d, err := subs.Next()
		if err != nil {
			break
		}
		_, _ = d.Ref.Update(ctx, []firestore.Update{
			{Path: "released", Value: true},
		})
	}

	return nil
}
func GetReleasedResult(examID, studentID string) (*models.Submission, error) {
	ctx := context.Background()
	client := config.DB

	doc := client.Collection("exams").Doc(examID).Collection("submissions").Doc(studentID)
	snap, err := doc.Get(ctx)
	if err != nil {
		return nil, err
	}

	var sub models.Submission
	if err := snap.DataTo(&sub); err != nil {
		return nil, err
	}

	if !sub.Released {
		return nil, errors.New("results not released yet")
	}

	return &sub, nil
}
func GetAllReleasedResultsForStudent(studentID string) ([]models.Submission, error) {
	ctx := context.Background()
	client := config.DB

	examsSnap, err := client.Collection("exams").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var results []models.Submission

	for _, examDoc := range examsSnap {
		examID := examDoc.Ref.ID

		subSnap, err := client.Collection("exams").
			Doc(examID).
			Collection("submissions").
			Doc(studentID).
			Get(ctx)

		if err != nil {
			continue
		}

		var sub models.Submission
		if err := subSnap.DataTo(&sub); err != nil {
			continue
		}

		if sub.Released {
			sub.StudentID = studentID
			sub.FinalScore = sub.AutoScore + sub.ManualScore
			results = append(results, sub)
		}
	}

	return results, nil
}
