package main

import (
	"context"
	"fmt"
	"github.com/teeverr/dummy_service/cmd/api/handlers"
	"github.com/teeverr/dummy_service/internal/workload"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/teeverr/dummy_service/internal/database"
	"github.com/teeverr/dummy_service/internal/domain"
	"github.com/teeverr/dummy_service/internal/health"

	"github.com/gorilla/mux"
)

func main() {
	config := &domain.Config{}
	err := config.Parse()
	if err != nil {
		log.Fatalf(fmt.Sprintf("can't parse configuration. Error: %s", err.Error()))
	}

	fmt.Println("connecting to database...")
	db, err := database.NewClient(config)
	if err != nil {
		log.Fatalf("something happened while db migrations came, %s", err.Error())
	}
	log.Printf("connection to the database is successfully established")

	router := mux.NewRouter()
	liveness := health.NewLivenessProbe()
	readiness := health.NewReadinessProbe()
	workloadHandlers := handlers.WorkloadHandler{Db: db, Config: config}
	checksTestHandlers := handlers.ChecksControl{Liveness: liveness, Readiness: readiness}

	router.Handle("/healthz", liveness)
	router.Handle("/readiness", readiness)
	router.Path("/healthz-on").Methods("GET").HandlerFunc(checksTestHandlers.HealthOn)
	router.Path("/healthz-off").Methods("GET").HandlerFunc(checksTestHandlers.HealthOff)
	router.Path("/ready").Methods("GET").HandlerFunc(checksTestHandlers.Ready)
	router.Path("/not-ready").Methods("GET").HandlerFunc(checksTestHandlers.NotReady)
	router.Path("/set-cpu").Methods("POST").HandlerFunc(workloadHandlers.SetCpuWorkload)
	router.Path("/show-cpu").Methods("GET").HandlerFunc(workloadHandlers.ShowTargetCPuWorkload)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		WriteTimeout: config.Server.WriteTimeout * time.Second,
		ReadTimeout:  config.Server.ReadTimeout * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	}()
	go workload.CpuWorkloadReader(db, config)

	readiness.Enable()
	log.Printf("web-server started")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	ctx, cancel := context.WithTimeout(context.Background(), config.Server.GracefulTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	readiness.Disable()
	fmt.Println("web-server shutting down")
}
