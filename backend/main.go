package main

import (
	"log"
	"os"

	"wsinspect/backend/common"
	"wsinspect/backend/core"
	"wsinspect/backend/routes"

	"github.com/spf13/cobra"
)

var (
	port   int
	target string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "wsinspect",
		Short: "WSInspect - WebSocket Debugging Tool",
		Long:  `A Postman-like debugging tool for WebSocket APIs`,
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the WSInspect proxy server",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize configuration
			config := core.NewConfig()
			if err := config.Load(); err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}

			// Initialize database
			db, err := core.NewDatabase(config)
			if err != nil {
				log.Fatalf("Failed to connect to database: %v", err)
			}

			// Auto migrate models
			if err := core.AutoMigrate(db); err != nil {
				log.Fatalf("Failed to migrate database: %v", err)
			}

			// Initialize services
			proxyService := core.NewProxyService(db)
			sessionService := core.NewSessionService(db)
			replayService := core.NewReplayService(db)
			fuzzService := core.NewFuzzService(db)

			// Setup router
			router := routes.SetupRouter(proxyService, sessionService, replayService, fuzzService)

			// Start server
			addr := ":" + config.GetString("server.port")
			log.Printf("Starting WSInspect server on %s", addr)
			if err := router.Run(addr); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		},
	}

	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	startCmd.Flags().StringVarP(&target, "target", "t", "ws://localhost:3000", "Target WebSocket server URL")

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
