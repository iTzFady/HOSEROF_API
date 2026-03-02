package models

import "time"

type LastSessionStats struct {
	SessionID  string    `json:"sessionId"`
	Date       time.Time `json:"date"`
	Attended   int       `json:"attended"`
	Absent     int       `json:"absent"`
	Percentage float64   `json:"percentage"`
}

type LowAttendanceStudent struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	AttendancePercentage float64 `json:"attendancePercentage"`
}

type TeacherDashboardResponse struct {
	TotalStudents         int                    `json:"totalStudents"`
	LastSession           LastSessionStats       `json:"lastSession"`
	Last4Sessions         []SessionChartItem     `json:"last4Sessions"`
	LowAttendanceStudents []LowAttendanceStudent `json:"lowAttendanceStudents"`
}
type SessionChartItem struct {
	Date     string `json:"date"`
	Attended int    `json:"attended"`
	Absent   int    `json:"absent"`
}

type AdminDashboardResponse struct {
	TotalStudents   int `json:"totalStudents"`
	TodayAttendance int `json:"todayAttendance"`
	PendingExams    int `json:"pendingExams"`
	PendingResults  int `json:"pendingResults"`
}
