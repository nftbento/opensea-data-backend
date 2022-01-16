/*

 */

package models

import (
	"time"
)

type AdminConfig struct {
	ID            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AdminUserID   string
	APIRate       uint
	OpenseaUrl    string
	OpenseaAPIKey string
}

func ReadAdminConfig(db *DB) (*AdminConfig, error) {
	config := AdminConfig{}
	err := db.Table("admin_configs").First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}
