package router

import (
	"context"
	"fmt"
	"gin-api/pkg/mysql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func NewHttpServer() {
	db := mysql.GetDB()

	r := gin.New()

	setAdminRouter(r, db)

	// pprof.Register(r)
	// r.Run(":8082")

	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown ")
	}
}
