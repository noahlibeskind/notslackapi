package clients

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

// Global Supabase Clients
var PublicClient *supabase.Client
var AuthClient *supabase.Client

func InitSupabase() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	supabaseURL := os.Getenv("SUPABASE_API_URL")
	serviceRoleKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	// Initialize Supabase clients
	PublicClient, err = supabase.NewClient(supabaseURL, serviceRoleKey, nil)
	if err != nil {
		log.Fatal("Failed to create public Supabase client:", err)
	}

	AuthClient, err = supabase.NewClient(supabaseURL, serviceRoleKey, nil)
	if err != nil {
		log.Fatal("Failed to create auth Supabase client:", err)
	}
}
