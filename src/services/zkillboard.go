package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tadeasf/eve-ran/src/db/models"
)

func FetchKillsFromZKillboard(characterID int64, page int) ([]models.Kill, error) {
	url := fmt.Sprintf("https://zkillboard.com/api/kills/characterID/%d/page/%d/", characterID, page)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "EVE Ran Application - GitHub: tadeasf/eve-ran")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawKills []struct {
		KillmailID int64 `json:"killmail_id"`
		ZKB        struct {
			LocationID     int64   `json:"locationID"`
			Hash           string  `json:"hash"`
			FittedValue    float64 `json:"fittedValue"`
			DroppedValue   float64 `json:"droppedValue"`
			DestroyedValue float64 `json:"destroyedValue"`
			TotalValue     float64 `json:"totalValue"`
			Points         int     `json:"points"`
			NPC            bool    `json:"npc"`
			Solo           bool    `json:"solo"`
			Awox           bool    `json:"awox"`
		} `json:"zkb"`
	}

	err = json.NewDecoder(resp.Body).Decode(&rawKills)
	if err != nil {
		return nil, err
	}

	var kills []models.Kill
	for _, rawKill := range rawKills {
		kill := models.Kill{
			KillmailID:     rawKill.KillmailID,
			CharacterID:    characterID,
			LocationID:     rawKill.ZKB.LocationID,
			Hash:           rawKill.ZKB.Hash,
			FittedValue:    rawKill.ZKB.FittedValue,
			DroppedValue:   rawKill.ZKB.DroppedValue,
			DestroyedValue: rawKill.ZKB.DestroyedValue,
			TotalValue:     rawKill.ZKB.TotalValue,
			Points:         rawKill.ZKB.Points,
			NPC:            rawKill.ZKB.NPC,
			Solo:           rawKill.ZKB.Solo,
			Awox:           rawKill.ZKB.Awox,
		}
		kills = append(kills, kill)
	}

	return kills, nil
}
