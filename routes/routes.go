package routes

import (
	"HOSEROF_API/controllers"
	"HOSEROF_API/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())
	protected.GET("/loginWithToken", controllers.TokenCheck)

	attendance := r.Group("/attendance")
	attendance.Use(middleware.RequireAuth())

	attendanceAdmin := attendance.Group("/")
	attendanceAdmin.Use(middleware.RequireAdmin())

	{
		attendanceAdmin.POST("/mark", controllers.MarkAttendance)
		attendanceAdmin.GET("/get/:studentID", controllers.GetAttendanceByID)
	}

	{
		attendance.GET("/get", controllers.GetAttendance)
	}

	exam := r.Group("/exam")
	exam.Use(middleware.RequireAuth())

	exam.GET("/list", controllers.ListExamsForStudent)
	exam.GET("/:examID", controllers.GetExamForStudent)
	exam.POST("/submit/:examID", controllers.SubmitExam)

	examAdmin := exam.Group("/")
	examAdmin.Use(middleware.RequireAdmin())

	examAdmin.POST("/create", controllers.CreateExam)
	examAdmin.GET("/submissions/:examID", controllers.GetSubmissionsForExam)
	examAdmin.POST("/grade/:examID", controllers.GradeAnswer)
	examAdmin.POST("/release/:examID", controllers.ReleaseResultsHandler)

	curriculum := r.Group("/curriculum")
	curriculum.Use(middleware.RequireAdmin())
	curriculum.GET("/:id", controllers.GetCurriculumByID)
	curriculum.GET("/class/:class_id", controllers.GetCurriculumsByClass)

	curriculumAdmin := exam.Group("/")
	curriculumAdmin.Use(middleware.RequireAdmin())
	curriculumAdmin.POST("/upload", controllers.UpdateCurriculum)
	curriculumAdmin.PUT("/:id", controllers.UpdateCurriculum)
	curriculumAdmin.DELETE("/:id", controllers.DeleteCurriculum)

	return r
}
