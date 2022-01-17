package opensea

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/NFTActions/opensea-data-backend/models"
	"github.com/NFTActions/opensea-data-backend/services/config"
	"github.com/NFTActions/opensea-data-backend/utils/conversion"
	"github.com/sirupsen/logrus"
)

type OpenseaService struct {
	DB  *models.DB
	log *logrus.Entry

	conf *config.AdminConfig
}

func NewOpenseaService(db *models.DB, log *logrus.Logger, adminConfig *config.AdminConfig) *OpenseaService {
	return &OpenseaService{
		DB:  db,
		log: log.WithField("service", "opensea"),

		conf: adminConfig,
	}
}

func (osvc *OpenseaService) GetRecentOpenseaEvents() ([]models.Activity, error) {
	unixNow := time.Now().Unix()
	occurred_before := unixNow - unixNow%60
	occurred_after := occurred_before - 60
	requestURL := fmt.Sprintf("%s/events?event_type=successful&occurred_after=%d&occurred_before=%d&limit=100",
		osvc.conf.GetConfig().OpenseaUrl, occurred_after, occurred_before)
	log.Println(requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}
	req.Header.Set("X-API-KEY", osvc.conf.GetConfig().OpenseaAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error calling Opensea events API: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Opensea events API returns error status code - %d", resp.StatusCode)
	} else {
		defer resp.Body.Close()
	}

	var responseBody ResponseBodyOpenseaEvents
	bodyInBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyInBytes, &responseBody)
	if err != nil {
		return nil, fmt.Errorf("Error in unmarshalling json payload from Opensea events API: %v", err)
	}

	recentActivities := []models.Activity{}
	for _, ae := range responseBody.AssetEvents {
		var createdDateParsed, listingTimeParsed, txTimestampParsed time.Time
		if ae.Asset.Collection.CreatedDate != "" {
			createdDateParsed, err = time.Parse("2006-01-02T15:04:05.000000", ae.Asset.Collection.CreatedDate)
			if err != nil {
				return nil, fmt.Errorf("Error in Parse collection created time - %s: %v", ae.Asset.Collection.CreatedDate, err)
			}
		}
		if ae.ListingTime != "" {
			listingTimeParsed, err = time.Parse("2006-01-02T15:04:05", ae.ListingTime)
			if err != nil {
				listingTimeParsed, err = time.Parse("2006-01-02T15:04:05.000000", ae.ListingTime)
				if err != nil {
					return nil, fmt.Errorf("Error in Parse listing time - %s: %v", ae.ListingTime, err)
				}
			}
		}
		txTimestampParsed, err = time.Parse("2006-01-02T15:04:05", ae.Transaction.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("Error in Parse txTimestamp - %s: %v", ae.Transaction.Timestamp, err)
		}
		activity := models.Activity{
			ID:                    ae.ID,
			TotalPrice:            conversion.ConvertWeiToGwei(ae.TotalPrice),
			CollectionCreatedDate: createdDateParsed,
			ListingTime:           listingTimeParsed,
			TxTimestamp:           txTimestampParsed,
			AssetName:             ae.Asset.Name,
			CollectionImageUrl:    ae.Asset.Collection.ImageUrl,
			CollectionName:        ae.Asset.Collection.Name,
			CollectionSlug:        ae.CollectionSlug,
			ContractAddress:       ae.Asset.AssetContract.Address,
			EventType:             ae.EventType,
			SellerAddress:         ae.Seller.Address,
			TransactionHash:       ae.Transaction.TransactionHash,
			WinnerAddress:         ae.WinnerAccount.Address,
		}
		recentActivities = append(recentActivities, activity)
	}
	return recentActivities, nil
}

type ResponseBodyOpenseaEvents struct {
	AssetEvents []AssetEvent `json:"asset_events"`
}

type AssetEvent struct {
	Asset          Asset       `json:"asset"`
	Seller         Seller      `json:"seller"`
	Transaction    Transaction `json:"transaction"`
	WinnerAccount  Winner      `json:"winner_account"`
	CollectionSlug string      `json:"collection_slug"`
	EventType      string      `json:"event_type"`
	TotalPrice     string      `json:"total_price"`
	ID             int64       `json:"id"`
	ListingTime    string      `json:"listing_time"`
}

type Asset struct {
	AssetContract AssetContract `json:"asset_contract"`
	Collection    Collection    `json:"collection"`
	Name          string        `json:"name"`
}

type AssetContract struct {
	Address string `json:"address"`
}

type Collection struct {
	CreatedDate string `json:"created_date"`
	ImageUrl    string `json:"image_url"`
	Name        string `json:"name"`
}

type Seller struct {
	Address string `json:"address"`
}

type Transaction struct {
	Timestamp       string `json:"timestamp"`
	TransactionHash string `json:"transaction_hash"`
}

type Winner struct {
	Address string `json:"address"`
}
