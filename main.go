// Application which greets you.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rich7690/promcord/internal/discord"
	"github.com/rich7690/promcord/internal/server"
)

var (
	Commit  string
	Version string
)

var (
	DiscordToken   string = os.Getenv("DISCORD_TOKEN")
	PerspectiveKey string = os.Getenv("PERSPECTIVE_KEY")
)

func main() {
	log.Println("Starting Server")
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, os.Interrupt)

	go func() {
		err := server.StartServer(ctx)
		if err != nil && err == http.ErrServerClosed {
			log.Printf("Error: %v\n", err)
		}
	}()

	d, err := discord.StartBot(ctx, DiscordToken, PerspectiveKey)
	if err != nil {
		log.Fatal(err)
	}
	err = d.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	<-sigs
	cancel()
	log.Println("Exiting")
	os.Exit(0)
}
