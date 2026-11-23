package main

import (
	"context"
	"fmt"
	"microservices/gocourse_user/internal/user"
	"microservices/gocourse_user/pkg/bootstrap"
	"microservices/gocourse_user/pkg/handler"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)
	logger := bootstrap.InitLogger()

	db, err := bootstrap.DBConnector()
	if err != nil {
		logger.Fatal("Failed to connect to the database: ", err)
	}
	pageLimitDefault := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pageLimitDefault == "" {
		logger.Fatal("PAGINATOR_LIMIT_DEFAULT environment variable is not set")
	}

	ctx := context.Background()
	userRepo := user.NewRepository(logger, db)
	userSrv := user.NewService(logger, userRepo)
	router := handler.NewUserHTTPServer(ctx, userSrv, user.Config{LimitPageDefault: pageLimitDefault})

	srv := &http.Server{
		Handler:      accessControlMiddleware(router),
		Addr:         address,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Println("Starting server on", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		logger.Fatal("Server error: ", err)
	}
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
