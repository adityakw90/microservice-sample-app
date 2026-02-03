package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/client"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/config"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/handler"
	httpMiddleware "github.com/adityakw90/microservice-sample-app/api-gateway/pkg/http"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Configuration: %s", cfg.String())

	// Initialize user gRPC client
	userClient, err := client.NewUserClient(client.Config{
		ServiceAddress: cfg.UserServiceAddress,
		DialTimeout:    cfg.UserServiceTimeout,
	})
	if err != nil {
		log.Fatalf("Failed to create user client: %v", err)
	}
	defer userClient.Close()

	// Initialize services
	userSvc := service.NewUserService(userClient.UserService())
	authSvc := service.NewAuthService(userClient.AuthService())

	// Initialize handlers
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(authSvc)

	// Create router
	router := mux.NewRouter()

	// Register API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	userHandler.RegisterRoutes(api)
	authHandler.RegisterRoutes(api)

	// Health check endpoint (no auth required)
	router.HandleFunc("/health", healthCheckHandler).Methods(http.MethodGet)

	// Apply middleware (wrap the router)
	middlewareChain := httpMiddleware.Chain(
		httpMiddleware.RecoveryMiddleware,
		httpMiddleware.CORSMiddleware,
	)(router)

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      middlewareChain,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("API Gateway starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// healthCheckHandler returns the health status of the API Gateway
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","service":"api-gateway"}`)
}
