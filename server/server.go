/*

 */

package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NFTActions/opensea-data-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	db         *models.DB
	log        *logrus.Logger
	router     *gin.Engine
	service    *Service
	controller *Controller
}

func CreateServer() *http.Server {

	/*
	 configure Logger
	*/
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = 4 // Info
	log.Formatter = &logrus.JSONFormatter{}

	/*
	 configure Database
	*/
	db_type, exists := os.LookupEnv("DB_TYPE")
	if !exists {
		log.Fatal("missing DB_TYPE environment variable")
	}
	db_path, exists := os.LookupEnv("DB_PATH")
	if !exists {
		log.Fatal("missing DB_PATH environment variable")
	}
	db := models.NewDB(db_type, db_path, log)
	if err := db.Connect(); err != nil {
		log.Fatal("db connection failed: ", err)
	}

	/*
		Initialize Server
	*/
	svr := NewServer(db, log)

	/*
		Initialize Services
	*/
	closers := svr.NewService()

	/*
		Initialize Controllers
	*/
	svr.controller = svr.NewController()

	/*
		Initialize Router
	*/
	svr.router = NewRouter(svr)

	/*
		Start HTTP Server
	*/
	// initialize server
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8080)
	httpServer := makeHttpServer(addr, svr.router)

	//todo: make socket server available for notification

	// handle graceful shutdown
	go handleGracefulShutdown(httpServer, closers)

	return httpServer
}

func NewServer(db *models.DB, log *logrus.Logger) *Server {
	return &Server{
		db:  db,
		log: log,
	}
}

func Start() error {
	srv := CreateServer()

	// listen and serve
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Println("server shutting down gracefully...")
	} else {
		log.Println("unexpected server shutdown...")
		log.Println("ERR: ", err)
	}
	return err
}
