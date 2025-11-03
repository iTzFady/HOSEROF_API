package main

import (
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
)

var studentsInFirestore = client.Collection("students")

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func loadStudentDocument(s string) *firestore.DocumentRef {
	return studentsInFirestore.Doc(s)
}
func studentLoginToFirestore(user User) (*UserData, error) {
	fmt.Println("function called")
	studentDocument := loadStudentDocument(user.StudentId)
	if studentDocument == nil {
		fmt.Println("user not found")
		return nil, errors.New("user not found")
	}
	studentDocumentAsMap, err := studentDocument.Get(ctx)
	if err != nil {
		fmt.Println("failed to map user")
		return nil, errors.New("failed to map user")
	}
	var userInFirestore UserInFirestore
	if err := studentDocumentAsMap.DataTo(&userInFirestore); err != nil {
		fmt.Println("failed to assign user to struct")
		return nil, errors.New("failed to assign user to struct")
	}
	var (
		studentToken string
		err1         error
		userData     UserData
	)
	fmt.Println(user.StudentPassword, userInFirestore.FirestoreStudentPassword)
	if user.StudentId == userInFirestore.FirestoreStudentId && CheckPasswordHash(user.StudentPassword, userInFirestore.FirestoreStudentPassword) {
		userData = mapUserDataToResponse(userInFirestore)
		fmt.Println("user verified")
		studentToken, err1 = jwtGenerator(user.StudentId, userData.StudentClass, userData.StudentName)
		if err1 != nil {
			fmt.Println("failed to return token")
			return nil, errors.New("failed to return token")
		}
	} else {
		fmt.Println("user access denied. invalid id/password")
	}
	fmt.Println("successfully performed operations")
	userData.StudentToken = studentToken
	return &userData, nil
}
func mapUserDataToResponse(f UserInFirestore) UserData {
	return UserData{
		StudentId:    f.FirestoreStudentId,
		StudentName:  f.FirestoreStudentName,
		StudentClass: f.FirestoreStudentClass,
	}

}
