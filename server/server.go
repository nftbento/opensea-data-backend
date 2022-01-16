/*

 */

package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NFTActions/opensea-data-backend/config"
	"github.com/NFTActions/opensea-data-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config     *config.Config
	db         *models.DB
	log        *logrus.Logger
	router     *gin.Engine
	service    *Service
	controller *Controller
}

func CreateServer() *http.Server {

	conf, err := config.NewConfig()
	if err != nil {
		panic("error reading config, " + err.Error())
	}

	log := logrus.New()
	log.Out = os.Stdout
	log.Level = conf.LogLevel()

	if conf.LogFileLocation() == "" {
		log.Fatal("missing log_file_location config variable")
	}
	logfile, err := os.OpenFile(conf.LogFileLocation(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("failed to open file for logging")
	} else {
		log.Out = logfile
		log.Formatter = &logrus.JSONFormatter{}
	}

	/*
	 configure Database
	*/
	cred := conf.DBCredentials()
	db := models.NewDB(cred[0], cred[1], log)

	if err := db.Connect(); err != nil {
		log.Fatal("db connection failed", err)
	}

	/*
		Initialize Server
	*/
	svr := NewServer(conf, db, log)

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
	svr.router = NewRouter(svr, *conf)

	/*
		Start HTTP Server
	*/
	// initialize server
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", conf.HTTPPort())
	httpServer := makeHttpServer(addr, svr.router)

	//todo: make socket server available for notification

	// handle graceful shutdown
	go handleGracefulShutdown(httpServer, closers)

	return httpServer
}

func NewServer(conf *config.Config, db *models.DB, log *logrus.Logger) *Server {
	return &Server{
		config: conf,
		db:     db,
		log:    log,
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
