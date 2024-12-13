package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}
var store = sessions.NewCookieStore([]byte("my-secret-key"))

func getSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "chat-session")
	return session
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Ambil file dari request body
		err := r.ParseMultipartForm(10 << 20) // Maksimal 10MB
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Ambil file dari form
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Membaca isi file
		var builder strings.Builder
		_, err = io.Copy(&builder, file)
		if err != nil {
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}

		// Proses file dan hitung konsumsi energi per ruangan
		energyConsumption, err := fileService.ProcessFile(builder.String())
		if err != nil {
			http.Error(w, "Failed to process file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Kirimkan hasil sebagai JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(energyConsumption)
	}).Methods("POST")

	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Context string `json:"context"`
			Query   string `json:"query"`
		}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if payload.Context == "" || payload.Query == "" {
			http.Error(w, "Context and query cannot be empty", http.StatusBadRequest)
			return
		}

		result, err := aiService.ChatWithAI(payload.Context, payload.Query, token)
		if err != nil {
			http.Error(w, "Failed to chat with AI: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
