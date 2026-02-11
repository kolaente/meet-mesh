// api/cmd/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kolaente/meet-mesh/api"
	gen "github.com/kolaente/meet-mesh/api/gen"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := api.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := api.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// Ensure avatar storage directory exists
	if err := os.MkdirAll(cfg.Storage.AvatarsPath, 0755); err != nil {
		log.Fatalf("Failed to create avatar storage directory: %v", err)
	}

	// Initialize auth service
	ctx := context.Background()
	auth, err := api.NewAuthService(ctx, &cfg.OIDC)
	if err != nil {
		log.Fatalf("Failed to init auth: %v", err)
	}

	// Initialize CalDAV client
	caldav := api.NewCalDAVClient(db)

	// Initialize mailer
	mailer, err := api.NewMailer(&cfg.SMTP, cfg.Server.BaseURL)
	if err != nil {
		log.Fatalf("Failed to init mailer: %v", err)
	}

	// Create handler
	handler := api.NewHandler(db, auth, caldav, mailer, cfg)

	// Create security handler
	security := api.NewSecurityHandler(db, auth)

	// Create avatar handler (plain HTTP, not ogen)
	avatarHandler := api.NewAvatarHandler(db, auth, cfg.Storage.AvatarsPath)

	// Create server with /api prefix so ogen strips it before routing
	server, err := gen.NewServer(handler, security, gen.WithPathPrefix("/api"))
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Create a mux that routes /api/* to the API server
	// and everything else to the static file server
	mux := http.NewServeMux()

	// Avatar routes (plain HTTP, not ogen - must be registered before /api/)
	mux.HandleFunc("/api/avatars/", avatarHandler.HandleAvatars)

	// API routes - the ogen server handles /api/*
	mux.Handle("/api/", server)

	// Static file server for everything else
	mux.Handle("/", api.NewStaticHandler())

	log.Printf("Meet Mesh starting on port %d...", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), mux))
}
