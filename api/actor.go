package api

import (
	"encoding/json"
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func AllActors(w http.ResponseWriter, r *http.Request) {
	logger.Log("Getting all actors...")

	colNCtx, err := db.GetCollectionAndContext(db.ActorCollection)
	if err != nil {
		logger.Error("allActors::Failed to get mongo actor collection")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cur, err := colNCtx.Col.Find(colNCtx.Ctx, bson.D{})
	if err != nil {
		logger.Warn("allActors::Failed to grab all actors", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	acts := make([]*models.Actor, 0)

	defer cur.Close(colNCtx.Ctx)

	for cur.Next(colNCtx.Ctx) {
		var a models.Actor
		err := cur.Decode(&a)
		if err != nil {
			logger.Warn("allActors::Failed to decode model with error", err)
			continue
		}

		acts = append(acts, &a)
	}

	err = json.NewEncoder(w).Encode(acts)
	if err != nil {
		logger.Error("allActors::Failed to encode actor json", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
