package config

import (
	"fmt"
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

var SupabaseStorage *storage_go.Client

func InitSupabase() error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		return fmt.Errorf("SUPABASE_URL and SUPABASE_SERVICE_KEY must be set")
	}

	SupabaseStorage = storage_go.NewClient(supabaseURL+"/storage/v1", supabaseKey, nil)

	return nil
}
