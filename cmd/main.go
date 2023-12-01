package main

import (
	"context"
	"os"
	"os/signal"

	conductor "firesidechuck.com/fat-controller/internal"
	"firesidechuck.com/fat-controller/internal/web"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config := conductor.LoadConfig()
	con := conductor.NewConductor(config.SerialPortName)

	go func() {
		con.Conduct(ctx)
	}()

	web.StartServer(ctx, con)
}
