package config

import (
	"errors"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

var (
	clientInstance *supabase.Client
	once           sync.Once
)

type SupabaseConfig struct {
	Url string
	Key string
}

func LoadConfig() (*SupabaseConfig, error) {
	_ = godotenv.Load(".env")

	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseUrl == "" || supabaseKey == "" {
		return nil, errors.New("Thiếu thông tin cấu hình Supabase")
	}

	return &SupabaseConfig{
		Url: supabaseUrl,
		Key: supabaseKey,
	}, nil
}

func InitSupabaseClient() *supabase.Client {
	once.Do(func() {
		config, err := LoadConfig()
		if err != nil {
			log.Fatalf("Lỗi khi tải cấu hình Supabase: %v", err)
		}

		client, err := supabase.NewClient(config.Url, config.Key, nil)

		if err != nil {
			log.Fatalf("Lỗi khi khởi tạo Supabase client: %v", err)
		}
		clientInstance = client
	})
	return clientInstance
}
