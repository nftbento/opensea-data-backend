/*

 */

package models

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	gorm.DB
	dbType           string
	dbConnectionPath string

	log *logrus.Entry
}

// NewDB initializes a DB object
func NewDB(dpType, dbConnectionPath string, log *logrus.Logger) *DB {
	return &DB{
		dbType:           dpType,
		dbConnectionPath: dbConnectionPath,
		log:              log.WithField("package", "models"),
	}
}

// Connect initiates a new connection with the given connection parameters
// for the database. It returns an error in case the connection fails.
func (db *DB) Connect() error {
	var dbConnection *gorm.DB
	var err error
	if db.dbType == "postgres" {
		dbConnection, err = gorm.Open(postgres.Open(db.dbConnectionPath), &gorm.Config{})
	} else {
		err = errors.New("invalid dbtype")
	}
	if err != nil {
		return err
	}
	_ = dbConnection.AutoMigrate(&Activity{})
	_ = dbConnection.AutoMigrate(&AdminConfig{})

	db.DB = *dbConnection

	return nil
}
