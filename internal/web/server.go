package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	conductor "firesidechuck.com/fat-controller/internal"
)

func StartServer(ctx context.Context, con *conductor.Conductor) {
	mux := http.NewServeMux()

	mux.HandleFunc("/active_train/horn", func(w http.ResponseWriter, r *http.Request) {
		con.Horn()
		status := http.StatusAccepted
		w.WriteHeader(status)
		fmt.Fprint(w, http.StatusText(status))
	})

	mux.HandleFunc("/active_train/lights/on", func(w http.ResponseWriter, r *http.Request) {
		con.LightsOn()
		status := http.StatusAccepted
		w.WriteHeader(status)
		fmt.Fprint(w, http.StatusText(status))
	})

	mux.HandleFunc("/active_train/lights/off", func(w http.ResponseWriter, r *http.Request) {
		con.LightsOff()
		status := http.StatusAccepted
		w.WriteHeader(status)
		fmt.Fprint(w, http.StatusText(status))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// Start a new context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
