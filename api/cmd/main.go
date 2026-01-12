// api/cmd/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

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

	// Create server
	server, err := gen.NewServer(handler, security)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Printf("Meet Mesh API starting on port %d...", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), server))
}
