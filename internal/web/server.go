package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	conductor "firesidechuck.com/fat-controller/internal"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartServer(
	ctx context.Context,
	log *slog.Logger,
	con *conductor.Conductor,
	adminToken string) {
	r := mux.NewRouter()

	r.HandleFunc("/active_train/horn", func(w http.ResponseWriter, r *http.Request) {
		con.Horn()
		status := http.StatusAccepted
		w.WriteHeader(status)
		fmt.Fprint(w, http.StatusText(status))
	})

	r.HandleFunc("/active_train/lights/on", func(w http.ResponseWriter, r *http.Request) {
		con.LightsOn()
		status := http.StatusAccepted
		w.WriteHeader(status)
		fmt.Fprint(w, http.StatusText(status))
	})

	r.HandleFunc("/active_train/lights/off", func(w http.ResponseWriter, r *http.Request) {
		con.LightsOff()
		status := http.StatusAccepted
		w.WriteHeader(status)
		fmt.Fprint(w, http.StatusText(status))
	})

	ar := r.PathPrefix("/admin").Subrouter()

	amw := NewAuthenticationMiddleware(adminToken)
	ar.Use(amw.Middleware)

	ar.HandleFunc("/track/on", func(w http.ResponseWriter, r *http.Request) {
		con.TrackPowerOn()
		status := http.StatusAccepted
		w.WriteHeader(status)
		fmt.Fprint(w, http.StatusText(status))
	})

	ar.HandleFunc("/track/off", func(w http.ResponseWriter, r *http.Request) {
		con.TrackPowerOff()
		status := http.StatusAccepted
		w.WriteHeader(status)
		fmt.Fprint(w, http.StatusText(status))
	})

	lr := handlers.CombinedLoggingHandler(os.Stdout, r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: lr,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	log.InfoContext(ctx, "shutting down gracefully, press Ctrl+C again to force")

	// Start a new context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.ErrorContext(ctx, "Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.InfoContext(ctx, "Server exiting")
}
