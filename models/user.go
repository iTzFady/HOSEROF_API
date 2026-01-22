package models

type UserLogin struct {
	StudentId string `json:"user_ID"`
}

type NewUser struct {
	NewStudentID          string `json:"student_id"`
	NewStudentName        string `json:"student_name"`
	NewStudentPhoneNumber string `json:"student_phonenumber"`
	NewStudentAge         string `json:"student_age"`
	NewStudentGrade       string `json:"student_grade"`
	NewStudentClass       string `json:"student_class"`
	NewStudentRole        string `json:"role"`
}

type UserFirestore struct {
	StudentID    string `firestore:"student_id"`
	StudentName  string `firestore:"student_name"`
	StudentClass string `firestore:"student_class"`
	StudentPhone string `firestore:"student_phonenumber"`
	StudentAge   string `firestore:"student_age"`
	StudentGrade string `firestore:"student_grade"`
	Role         string `firestore:"role"`
}

type UserClassList struct {
	StudentID    string `firestore:"student_id"`
	StudentName  string `firestore:"student_name"`
	StudentGrade string `firestore:"student_grade"`
}

type UserDataResponse struct {
	StudentToken string `json:"student_token"`
	StudentId    string `json:"student_id"`
	StudentName  string `json:"student_name"`
	StudentClass string `json:"student_class"`
	Role         string `json:"role"`
}
