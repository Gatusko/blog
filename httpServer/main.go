package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/Gatusko/blog/internal/data"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type resp struct {
	Message string `json:"error"`
}

type ready struct {
	Status string `json:"status"`
}

// For future Injection in our handlers
type apiConfig struct {
	db     string
	models data.Models
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Issue loading env file %s", err)
		return
	}
	port := os.Getenv("PORT")
	dbString := os.Getenv("DB")
	log.Printf("Running on port : %s", port)
	myApi := apiConfig{dbString, data.Models{}}
	db, err := myApi.openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//Injecting Database type of Connection
	myApi.models = data.NewModels(db)
	go myApi.workerScrappers(2, time.Second*5)
	log.Printf("Connection Worked")
	err = http.ListenAndServe(":"+port, myApi.mapAllRouters())
	log.Fatalf("Got an error running server: %s", err)
}

func (apiConfig *apiConfig) workerScrappers(number int, interval time.Duration) {
	ticker := time.NewTicker(interval)
	log.Printf("Enter of workScrapper")
	for range ticker.C {
		feeds, err := apiConfig.models.Feeds.GetNextFeedsToFetch(number)
		if err != nil {
			log.Printf("Error at getting the next Feeds : %s", err)
			continue
		}
		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)
			go apiConfig.ScrapData(feed.Id, &wg)
		}
		wg.Wait()
		log.Printf("Succes of scraping the data")
	}
}

func testingGo() {
	log.Println("test")
}

func (app *apiConfig) openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", app.db)
	if err != nil {
		return nil, err
	}
	// Make a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// If context doesn't respond in 5 minutes
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *apiConfig) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling json : %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("SErver error 5xx: %s", msg)
	}
	res := resp{msg}
	respondWithJSON(w, code, res)
}
