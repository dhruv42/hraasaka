package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dhruv42/hraasaka/cache"
	"github.com/dhruv42/hraasaka/config"
	"github.com/dhruv42/hraasaka/db"
	"github.com/dhruv42/hraasaka/enums"
	"github.com/dhruv42/hraasaka/helpers"
	"github.com/dhruv42/hraasaka/models"
)

type RequestShortenLink struct {
	Url string `json:"url,omitempty"`
}

type ResponseShortenLink struct {
	CreatedAt time.Time             `json:"createdAt,omitempty"`
	UpdatedAt time.Time             `json:"updatedAt,omitempty"`
	Url       string                `json:"url,omitempty"`
	Hash      string                `json:"hash,omitempty"`
	ID        interface{}           `json:"id,omitempty"`
	Status    enums.ShortLinkStatus `json:"status,omitempty"`
}

type response struct {
	Success bool        `json:"success"`
	Status  int         `json:"statusCode"`
	Data    interface{} `json:"data"`
}

const COLLECTION = "links"

func ShortenUrl(w http.ResponseWriter, r *http.Request) {

	var shortenReq RequestShortenLink

	ctx := context.Background()

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error loading config: [%v]", err)
	}

	client, err := db.GetDBConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("Db connectiom failed:- %v", err)
	}
	defer client.Disconnect(ctx)

	cache, err := cache.ConnectCache()
	if err != nil {
		log.Fatalf("Cache connectiom failed:- %v", err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		generateResponse(w, "BAD_REQUEST", http.StatusBadRequest, nil)
		return
	}

	err = json.Unmarshal(body, &shortenReq)
	if err != nil {
		generateResponse(w, "BAD_REQUEST", http.StatusBadRequest, nil)
		return
	}

	// Get counter from cache and update

	counter, err := cache.Get("counter").Result()
	if err != nil {
		log.Fatalf("Failed to set counter in cache:- %v", err)
	}

	counterInt, err := strconv.Atoi(counter)
	if err != nil {
		log.Fatalf("Failed to convert counter string into int:- %v", err)
	}

	hash := helpers.Base62Encode(counterInt)

	insertDoc := &models.Link{
		Url:       shortenReq.Url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    enums.ShortLinkCreated,
		Hash:      hash,
	}

	collection := client.Database(cfg.DbName).Collection(COLLECTION)
	insertRes, err := collection.InsertOne(ctx, insertDoc)
	if err != nil {
		generateResponse(w, "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, nil)
	}

	_ = cache.Incr("counter")

	finalResponse := &ResponseShortenLink{
		CreatedAt: insertDoc.CreatedAt,
		UpdatedAt: insertDoc.UpdatedAt,
		Status:    insertDoc.Status,
		Hash:      insertDoc.Hash,
		Url:       insertDoc.Url,
		ID:        insertRes.InsertedID,
	}

	generateResponse(w, "SUCCESS", http.StatusOK, finalResponse)

}

func RedirectUrl(w http.ResponseWriter, r *http.Request) {

}

func generateResponse(w http.ResponseWriter, message string, statusCode int, respBody interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// resp := make(map[string]string)
	resp := response{
		Success: true,
		Status:  statusCode,
		Data:    respBody,
	}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
