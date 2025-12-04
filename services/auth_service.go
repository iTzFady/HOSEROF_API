package services

import (
	"HOSEROF_API/config"
	"HOSEROF_API/models"
	"context"
	"errors"
	"fmt"
)

func SignupUser(newUser models.NewUser) error {
	data := map[string]interface{}{
		"student_id":          newUser.NewStudentID,
		"student_name":        newUser.NewStudentName,
		"student_phonenumber": newUser.NewStudentPhoneNumber,
		"student_age":         newUser.NewStudentAge,
		"student_grade":       newUser.NewStudentGrade,
		"student_class":       newUser.NewStudentClass,
		"role":                newUser.NewStudentRole,
	}

	_, err := config.DB.Collection("students").
		Doc(newUser.NewStudentID).
		Set(context.Background(), data)

	if err != nil {
		return err
	}

	return nil
}

func LoginUser(login models.UserLogin) (*models.UserDataResponse, error) {
	ctx := context.Background()
	docRef := config.DB.Collection("students").Doc(login.StudentId)
	docSnap, err := docRef.Get(ctx)

	var fsUser models.UserFirestore

	if err := docSnap.DataTo(&fsUser); err != nil {
		return nil, errors.New("INVALID_LOGIN_PAYLOAD")
	}

	token, err := jwtGenerator(fsUser.StudentID, fsUser.StudentClass, fsUser.Role, fsUser.StudentName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	resp := &models.UserDataResponse{
		StudentToken: token,
		StudentId:    fsUser.StudentID,
		StudentName:  fsUser.StudentName,
		StudentClass: fsUser.StudentClass,
		Role:         fsUser.Role,
	}
	return resp, nil
}
