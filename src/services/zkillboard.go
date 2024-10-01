package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawKills []struct {
		KillmailID    int64     `json:"killmail_id"`
		KillmailTime  time.Time `json:"killmail_time"`
		SolarSystemID int       `json:"solar_system_id"`
		ZKB           struct {
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
	err = json.Unmarshal(body, &rawKills)
	if err != nil {
		return nil, err
	}

	var kills []models.Kill
	for _, rawKill := range rawKills {
		kill := models.Kill{
			KillmailID:     rawKill.KillmailID,
			CharacterID:    characterID,
			KillTime:       rawKill.KillmailTime,
			SolarSystemID:  rawKill.SolarSystemID,
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

func FetchKillmailFromESI(killmailID int64, killmailHash string) (*models.Kill, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/killmails/%d/%s/", killmailID, killmailHash)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var esiKillmail struct {
		KillmailTime  string            `json:"killmail_time"`
		Victim        models.Victim     `json:"victim"`
		Attackers     []models.Attacker `json:"attackers"`
		SolarSystemID int               `json:"solar_system_id"`
	}
	err = json.Unmarshal(body, &esiKillmail)
	if err != nil {
		return nil, err
	}

	killTime, err := time.Parse(time.RFC3339, esiKillmail.KillmailTime)
	if err != nil {
		return nil, fmt.Errorf("error parsing kill time: %v", err)
	}

	return &models.Kill{
		KillmailID:    killmailID,
		KillTime:      killTime,
		Victim:        esiKillmail.Victim,
		Attackers:     esiKillmail.Attackers,
		SolarSystemID: esiKillmail.SolarSystemID,
	}, nil
}
