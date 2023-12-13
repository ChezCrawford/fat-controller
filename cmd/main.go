package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	conductor "firesidechuck.com/fat-controller/internal"
	"firesidechuck.com/fat-controller/internal/web"
)

func main() {
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	log := slog.New(h)
	slog.SetDefault(log)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config := conductor.LoadConfig()

	con := conductor.NewConductor(log, config.SerialPortName, config.UseSimDriver)

	go func() {
		con.Conduct(ctx)
	}()

	web.StartServer(ctx, log, con)
}
