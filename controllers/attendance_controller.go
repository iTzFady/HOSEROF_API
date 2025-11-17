package controllers

import (
	"HOSEROF_API/services"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

func MarkAttendance(c *gin.Context) {
	var body struct {
		StudentID string `json:"studentId"`
		Attended  bool   `json:"attended"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if body.StudentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentId required"})
		return
	}
	err := services.MarkAttendance(body.StudentID, body.Attended)

	if err != nil {

		if err.Error() == "no user found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save attendance"})
		log.Print(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})

}
func GetAttendance(c *gin.Context) {
	token := c.MustGet("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	studentID := claims["user_ID"].(string)

	resp, err := services.GetAttendance(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get attendance"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetAttendanceByID(c *gin.Context) {
	studentID := c.Param("studentID")

	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	resp, err := services.GetAttendance(studentID)

	if err != nil {

		if err.Error() == "no user found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get attendance"})
		return
	}
	c.JSON(http.StatusOK, resp)

}
