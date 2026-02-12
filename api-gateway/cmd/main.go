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

	grpcAdapter "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/secondary/grpc"
	userHandler "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/httpHandler"
	authHandler "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/httpHandler"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/application/user"
	authApp "github.com/adityakw90/microservice-sample-app/api-gateway/internal/application/auth"
	userFileApp "github.com/adityakw90/microservice-sample-app/api-gateway/internal/application/user_file"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/config"
	httpMiddleware "github.com/adityakw90/microservice-sample-app/api-gateway/pkg/http"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Configuration: %s", cfg.String())

	// === Secondary Adapters (Driven) ===
	// Create gRPC client adapters (implements both UserClient and AuthClient)
	log.Println("Initializing gRPC client adapters...")

	userGrpcAdapter, authGrpcAdapter, err := grpcAdapter.NewUserClientAdapter(cfg.UserServiceAddress)
	if err != nil {
		log.Fatalf("Failed to create user gRPC adapter: %v", err)
	}

	// === Application Services (Use Cases) ===
	// Create application services that depend on port interfaces
	log.Println("Initializing application services...")

	userAppService := user.NewUserApplicationService(userGrpcAdapter)
	authAppService := authApp.NewAuthApplicationService(authGrpcAdapter)

	// Get the underlying connection for UserFile gRPC adapter
	userGrpcAdapterWithConn := userGrpcAdapter.(*grpcAdapter.UserClientAdapter)
	userFileGrpcAdapter := grpcAdapter.NewUserFileClientAdapter(userGrpcAdapterWithConn.GetConn())
	userFileAppService := userFileApp.NewUserFileApplicationService(userFileGrpcAdapter)

	// === Primary Adapters (Driving) ===
	// Create HTTP handlers that use application services
	log.Println("Initializing HTTP handlers...")

	userHdlr := userHandler.NewUserHandler(userAppService)
	authHdlr := authHandler.NewAuthHandler(authAppService)
	userFileHdlr := authHandler.NewUserFileHandler(userFileAppService)

	// === Router Configuration ===
	// Create router and register routes
	log.Println("Registering routes...")

	router := mux.NewRouter()

	// Register API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	userHdlr.RegisterRoutes(api)
	authHdlr.RegisterRoutes(api)
	userFileHdlr.RegisterRoutes(api)

	// Health check endpoint (no auth required)
	router.HandleFunc("/health", healthCheckHandler).Methods(http.MethodGet)

	// Apply middleware (wrap the router)
	middlewareChain := httpMiddleware.Chain(
		httpMiddleware.RecoveryMiddleware,
		httpMiddleware.CORSMiddleware,
	)(router)

	// === Server Start ===
	// Create and start HTTP server
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
