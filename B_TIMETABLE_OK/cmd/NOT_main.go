package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"

	"middleware/example/internal/controllers"
	"middleware/example/internal/repositories"
	"middleware/example/internal/services"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./users.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	// important pour faire respecter la FK alerts.agenda_id -> agendas.id
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	agendaRepo := repositories.NewAgendaSQLite(db)
	if err := agendaRepo.Init(context.Background()); err != nil {
		log.Fatal(err)
	}
	agendaSvc := services.NewAgendaService(agendaRepo)
	agendaCtrl := controllers.NewAgendaController(agendaSvc)

	// --- Alerts wiring ---
	alertRepo := repositories.NewAlertSQLite(db)
	if err := alertRepo.Init(context.Background()); err != nil {
		log.Fatal(err)
	}
	alertSvc := services.NewAlertService(alertRepo, agendaRepo)
	alertCtrl := controllers.NewAlertController(alertSvc)

	r := chi.NewRouter()

	r.Handle("/api/*", http.StripPrefix("/api", http.FileServer(http.Dir("./api"))))

	// UI Swagger accessible sur http://localhost:8080/swagger/index.html
r.Get("/swagger/*", httpSwagger.Handler(
    httpSwagger.URL("/api/swagger.json"), // l’UI pointe vers ce JSON
    ))



	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Mount("/agendas", agendaCtrl.Routes())
	r.Mount("/alerts", alertCtrl.Routes())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("listening on :" + port)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
