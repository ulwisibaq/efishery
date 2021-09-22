package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/ulwisibaq/efishery/auth/internal/auth/repository"
	"github.com/ulwisibaq/efishery/auth/internal/auth/service"
	"github.com/ulwisibaq/efishery/auth/internal/handler"
	"github.com/ulwisibaq/efishery/config"
)

var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Starts REST API",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.MainConfig{}
		config.ReadConfig(cfg)

		// setup mysql db conn
		db, err := sqlx.Open(`mysql`, cfg.Database.MysqlDSN)
		if err != nil {
			log.Fatalf("couldn't connect to mysql db, error :%v", err)
			return err
		}

		err = db.Ping()
		if err != nil {
			log.Fatal(err)
			return err
		}

		defer func() {
			err := db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		authRepository := repository.NewAuthRepository(db)
		authService := service.NewAuthService(authRepository)
		authHandler := handler.NewAuthHandler(authService)

		myRouter := mux.NewRouter()
		myRouter.HandleFunc("/user/register", authHandler.RegisterUser).Methods("POST")
		myRouter.HandleFunc("/user/login", authHandler.Login).Methods("POST")
		myRouter.HandleFunc("/user/verify", authHandler.Verify).Methods("POST")

		srv := http.Server{
			Addr:    ":8081",
			Handler: myRouter,
		}
		go func() {
			log.Println("auth http server starting in port :8081")
			if err := srv.ListenAndServe(); err != nil {
				log.Fatal(err)
			}
		}()

		// graceful shutdown
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		log.Println("server stopped")

		return err
	},
}
