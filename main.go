package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/RodolfoCamposGlz/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = "8080"
const filepathRoot = "."

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// Return an http.Handler that wraps the next handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment your metric
		cfg.fileserverHits.Add(1)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits.Load())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func (cfg *apiConfig) handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Not allowed")
		return
	}
	err := cfg.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users")
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}



var handler http.Handler = http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

func main (){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	fmt.Println("DB_URL: ", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	apiCfg.dbQueries = dbQueries
	apiCfg.platform = platform
	mux := http.NewServeMux()

	// Create a new http.Server
	server := &http.Server{
		Addr:  ":" + port, // Bind to port 8080
		Handler: mux,     // Use the ServeMux as the handler
	}

	mux.Handle("/app/",apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", apiCfg.handlerReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Start the server
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}



