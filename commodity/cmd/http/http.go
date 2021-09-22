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
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cobra"
	"github.com/ulwisibaq/efishery/commodity/internal/commodity/repository"
	"github.com/ulwisibaq/efishery/commodity/internal/commodity/service"
	"github.com/ulwisibaq/efishery/commodity/internal/handler"
	"github.com/ulwisibaq/efishery/config"
)

var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Starts REST API",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cfg := &config.MainConfig{}
		config.ReadConfig(cfg)

		cacheCommodity := cache.New(23*time.Hour, 24*time.Hour)

		commodityRepository := repository.NewCommodityRepository(cfg.CommoditiesAPI.Url)
		commodityService := service.NewCommodityService(commodityRepository, cacheCommodity)
		commodityHandler := handler.NewCommodityHandler(commodityService)

		myRouter := mux.NewRouter()
		myRouter.HandleFunc("/commodity/fetch", commodityHandler.GetCommodity).Methods("GET")
		myRouter.HandleFunc("/commodity/verify", commodityHandler.Verify).Methods("POST")

		srv := http.Server{
			Addr:    ":8082",
			Handler: myRouter,
		}
		go func() {
			log.Println("commodity http server starting in port :8082")
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
