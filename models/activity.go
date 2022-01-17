package models

import (
	"time"
)

type Activity struct {
	ID                    int64     `json:"id"  gorm:"primaryKey"`
	TotalPrice            int64     `json:"total_price"`
	CollectionCreatedDate time.Time `json:"collection_created_date"`
	ListingTime           time.Time `json:"listing_time"`
	TxTimestamp           time.Time `json:"tx_timestamp"`
	AssetName             string    `json:"asset_name"`
	CollectionImageUrl    string    `json:"collection_image_url"`
	CollectionName        string    `json:"collection_name"`
	CollectionSlug        string    `json:"collection_slug"`
	ContractAddress       string    `json:"contract_address"`
	EventType             string    `json:"event_type"`
	SellerAddress         string    `json:"seller_address"`
	TransactionHash       string    `json:"transaction_hash"`
	WinnerAddress         string    `json:"winner_address"`
}

func (db *DB) InsertActivity(activity *Activity) error {
	return db.Create(activity).Error
}

func (db *DB) BatchInsertActivity(activities []Activity) error {
	return db.Create(activities).Error
}

func (db *DB) GetActivitiesAfter(occurred_after time.Time) ([]Activity, error) {
	activities := []Activity{}
	err := db.Where("tx_timestamp > ?", occurred_after).Find(&activities).Error
	if err != nil {
		return nil, err
	}
	return activities, nil
}
