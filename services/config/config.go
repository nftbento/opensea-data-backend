package config

import (
	"github.com/sirupsen/logrus"
	"opensea-data-backend/models"
)

type AdminConfig struct {
	DB  *models.DB
	log *logrus.Logger

	conf *models.AdminConfig
	quit chan int
}

func NewAdminConfig(db *models.DB, logger *logrus.Logger) *AdminConfig {
	s := &AdminConfig{
		DB:  db,
		log: logger,
	}

	s.updateConf()

	// ticker := time.NewTicker(5 * time.Minute)
	// quit := make(chan int)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			s.updateConf()
	// 		case <-quit:
	// 			ticker.Stop()
	// 			return
	// 		}
	// 	}
	// }()
	// s.quit = quit
	return s
}

func (s *AdminConfig) updateConf() {
	adminConfig, err := models.ReadAdminConfig(s.DB)

	if err != nil {
		s.log.Error("failed to fetch latest config")
	}
	s.conf = adminConfig
}

func (s *AdminConfig) GetConfig() *models.AdminConfig {
	return s.conf
}

func (s *AdminConfig) Close() {
	s.quit <- 0
}
