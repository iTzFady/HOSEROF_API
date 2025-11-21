package services

import (
	"HOSEROF_API/config"
	"HOSEROF_API/models"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	storage_go "github.com/supabase-community/storage-go"
	"google.golang.org/api/iterator"
)

var supabaseStorage *storage_go.Client

func InitSupabaseStorage() error {
	supabaseStorage = config.SupabaseStorage
	return nil
}

func UploadCurriculum(ctx context.Context, req models.UploadCurriculumRequest, file multipart.File, header *multipart.FileHeader, userID string) (*models.Curriculum, error) {
	ext := filepath.Ext(header.Filename)
	uniqueFilename := fmt.Sprintf("%s_%s%s", req.ClassID, uuid.New().String(), ext)
	storagePath := fmt.Sprintf("curriculum/%s/%s", req.ClassID, uniqueFilename)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	fileReader := bytes.NewReader(fileBytes)

	_, err = config.SupabaseStorage.UploadFile("curriculum", storagePath, fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to supabase: %w", err)
	}

	resp := config.SupabaseStorage.GetPublicUrl("curriculum", storagePath)
	fileURL := resp.SignedURL

	fileType := getFileType(ext)

	curriculum := &models.Curriculum{
		ID:          uuid.New().String(),
		ClassID:     req.ClassID,
		Title:       req.Title,
		Description: req.Description,
		FileType:    fileType,
		FileURL:     fileURL,
		FileName:    header.Filename,
		FileSize:    header.Size,
		Order:       req.Order,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = config.DB.Collection("curriculums").Doc(curriculum.ID).Set(ctx, curriculum)
	if err != nil {
		config.SupabaseStorage.RemoveFile("curriculum", []string{storagePath})
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	return curriculum, nil
}

func GetCurriculumsByClass(ctx context.Context, classID string) ([]models.Curriculum, error) {
	var curriculums []models.Curriculum

	iter := config.DB.Collection("curriculums").
		Where("class_id", "==", classID).
		OrderBy("order", firestore.Asc).
		OrderBy("created_at", firestore.Asc).
		Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate documents: %w", err)
		}

		var curriculum models.Curriculum
		if err := doc.DataTo(&curriculum); err != nil {
			return nil, fmt.Errorf("failed to parse curriculum: %w", err)
		}
		curriculums = append(curriculums, curriculum)
	}

	return curriculums, nil
}

func GetAllCurriculums(ctx context.Context) ([]models.Curriculum, error) {
	var curriculums []models.Curriculum

	iter := config.DB.Collection("curriculums").
		OrderBy("class_id", firestore.Asc).
		OrderBy("order", firestore.Asc).
		Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate documents: %w", err)
		}

		var curriculum models.Curriculum
		if err := doc.DataTo(&curriculum); err != nil {
			return nil, fmt.Errorf("failed to parse curriculum: %w", err)
		}
		curriculums = append(curriculums, curriculum)
	}

	return curriculums, nil
}

func GetCurriculumByID(ctx context.Context, id string) (*models.Curriculum, error) {
	doc, err := config.DB.Collection("curriculums").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("curriculum not found: %w", err)
	}

	var curriculum models.Curriculum
	if err := doc.DataTo(&curriculum); err != nil {
		return nil, fmt.Errorf("failed to parse curriculum: %w", err)
	}

	return &curriculum, nil
}

func UpdateCurriculum(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	_, err := config.DB.Collection("curriculums").Doc(id).Update(ctx, []firestore.Update{
		{Path: "title", Value: updates["title"]},
		{Path: "description", Value: updates["description"]},
		{Path: "order", Value: updates["order"]},
		{Path: "updated_at", Value: updates["updated_at"]},
	})

	return err
}

func DeleteCurriculum(ctx context.Context, id string) error {
	curriculum, err := GetCurriculumByID(ctx, id)
	if err != nil {
		return err
	}

	ext := filepath.Ext(curriculum.FileName)
	filenameFromID := fmt.Sprintf("%s_%s%s", curriculum.ClassID, id, ext)
	storagePath := fmt.Sprintf("curriculum/%s/%s", curriculum.ClassID, filenameFromID)

	_, err = config.SupabaseStorage.RemoveFile("curriculum", []string{storagePath})
	if err != nil {
		fmt.Printf("Warning: failed to delete file from Supabase: %v\n", err)
	}

	_, err = config.DB.Collection("curriculums").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete from firestore: %w", err)
	}

	return nil
}

func getFileType(ext string) string {
	switch ext {
	case ".pdf":
		return "pdf"
	case ".mp3", ".wav", ".ogg", ".m4a":
		return "audio"
	case ".mp4", ".avi", ".mov", ".wmv":
		return "video"
	case ".doc", ".docx", ".txt", ".rtf":
		return "document"
	case ".ppt", ".pptx":
		return "presentation"
	case ".xls", ".xlsx":
		return "spreadsheet"
	case ".jpg", ".jpeg", ".png", ".gif":
		return "image"
	default:
		return "other"
	}
}
