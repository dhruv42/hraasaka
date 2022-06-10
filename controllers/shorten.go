package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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

	hash := helpers.Base62Encode(int(helpers.GetCounter()))

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
