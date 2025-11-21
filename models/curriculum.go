package models

import "time"

type Curriculum struct {
	ID          string    `json:"id" firestore:"id"`
	ClassID     string    `json:"class_id" firestore:"class_id"`
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	FileType    string    `json:"file_type" firestore:"file_type"`
	FileURL     string    `json:"file_url" firestore:"file_url"`
	FileName    string    `json:"file_name" firestore:"file_name"`
	FileSize    int64     `json:"file_size" firestore:"file_size"`
	Order       int       `json:"order" firestore:"order"`
	CreatedBy   string    `json:"created_by" firestore:"created_by"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" firestore:"updated_at"`
}

type UploadCurriculumRequest struct {
	ClassID     string `json:"class_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type CurriculumListResponse struct {
	Curriculums []Curriculum `json:"curriculums"`
	Total       int          `json:"total"`
}
