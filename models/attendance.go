package models

import "time"

type AttendanceRecord struct {
	Attended  bool      `json:"attended" firestore:"attended"`
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}

type AttendanceResponse struct {
	Records    []AttendanceRecord `json:"records"`
	Percentage float64            `json:"percentage"`
}

type MarkAttendanceRequest struct {
	Attended bool `json:"attended"`
}
