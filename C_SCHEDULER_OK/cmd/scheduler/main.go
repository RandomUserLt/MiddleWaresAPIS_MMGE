package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	zsched "github.com/zhashkevych/scheduler"

	"middleware/scheduler/internal/configclient"
	"middleware/scheduler/internal/natsbus"
	"middleware/scheduler/internal/scheduler"
	"middleware/scheduler/internal/timetableclient"
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	cfgURL := env("CONFIG_URL", "http://localhost:8080")
	ttURL := env("TIMETABLE_URL", "http://localhost:8081")
	natsURL := env("NATS_URL", "nats://127.0.0.1:4222")
	period := env("SCHEDULER_PERIOD", "10m") // ex: "5m", "1h", etc.

	// --- initialisation des clients ---
	cfg := configclient.New(cfgURL)
	tt := timetableclient.New(ttURL)
	pub, err := natsbus.New(context.Background(), natsbus.Options{
		URL:     natsURL,
		Stream:  "EVENTS",
		Subject: "EVENTS.new",
	})
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer pub.Close()

	job := scheduler.NewJob(cfg, tt, pub)

	// Optionnel : réduire la période temporelle
	// job.From = "2025-11-01"
	// job.To   = "2025-11-30"

	// --- mode exécution unique ---
	if env("RUN_ONCE", "") != "" {
		if err := job.RunOnce(context.Background()); err != nil {
			log.Fatalf("[scheduler] error: %v", err)
		}
		return
	}

	// --- planification périodique ---
	dur, err := time.ParseDuration(period)
	if err != nil {
		dur = 10 * time.Minute
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sc := zsched.NewScheduler()

	// Ajoute le job (le scheduler tourne automatiquement en arrière-plan)
	sc.Add(ctx, func(c context.Context) {
		if err := job.RunOnce(c); err != nil {
			log.Printf("[scheduler] run error: %v", err)
		}
	}, dur)

	log.Printf("[scheduler] running every %s (using zhashkevych/scheduler)", dur)

	// Bloc d'attente pour arrêt propre
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	cancel()
	log.Println("[scheduler] stopped")
}
