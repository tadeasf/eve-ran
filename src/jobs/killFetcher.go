package jobs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/services"
)

func StartKillFetcherJob() {
	c := cron.New()
	c.AddFunc("@every 1h", func() {
		log.Println("Starting to fetch kills for all characters")
		fetchKillsForAllCharacters()
	})
	c.Start()

	// Run the job immediately on startup
	go fetchKillsForAllCharacters()
}

func fetchKillsForAllCharacters() {
	characters, err := db.GetAllCharacters()
	if err != nil {
		log.Printf("Error fetching characters: %v", err)
		return
	}

	log.Printf("Found %d characters", len(characters))

	for _, character := range characters {
		fetchKillsForCharacter(character.ID)
	}

	log.Println("Finished fetching kills for all characters")
}

func fetchKillsForCharacter(characterID int64) {
	lastKillTime, err := db.GetLastKillTimeForCharacter(characterID)
	if err != nil {
		log.Printf("Error getting last kill time for character %d: %v", characterID, err)
		lastKillTime = time.Time{}
	}
	log.Printf("Last kill time for character %d: %v", characterID, lastKillTime)

	isNewCharacter := lastKillTime.IsZero()
	page := 1
	totalNewKills := 0

	const maxConcurrentRequests = 200
	semaphore := make(chan struct{}, maxConcurrentRequests)
	var wg sync.WaitGroup

	for {
		log.Printf("Fetching page %d for character %d", page, characterID)
		kills, err := services.FetchKillsFromZKillboard(characterID, page)
		if err != nil {
			log.Printf("Error fetching kills for character %d: %v", characterID, err)
			return
		}

		if len(kills) == 0 {
			log.Printf("No more kills found for character %d", characterID)
			break
		}

		newKills := 0
		killChan := make(chan *models.Kill, len(kills))

		for _, kill := range kills {
			if isNewCharacter || kill.KillTime.After(lastKillTime) {
				wg.Add(1)
				go func(k models.Kill) {
					defer wg.Done()
					semaphore <- struct{}{}
					defer func() { <-semaphore }()

					esiKill, err := services.FetchKillmailFromESI(k.KillmailID, k.Hash)
					if err != nil {
						log.Printf("Error fetching ESI killmail %d: %v", k.KillmailID, err)
						return
					}

					// Combine zKillboard and ESI data
					k.Victim = esiKill.Victim
					k.Attackers = esiKill.Attackers
					killChan <- &k
				}(kill)
				newKills++
			} else if !isNewCharacter {
				log.Printf("Reached already processed kills for character %d", characterID)
				close(killChan)
				wg.Wait()
				return
			}
		}

		go func() {
			wg.Wait()
			close(killChan)
		}()

		for kill := range killChan {
			err = db.InsertKill(kill)
			if err != nil {
				log.Printf("Error inserting kill %d: %v", kill.KillmailID, err)
			} else {
				totalNewKills++
			}
		}

		log.Printf("Inserted %d new kills for character %d on page %d", newKills, characterID, page)

		if newKills == 0 && !isNewCharacter {
			log.Printf("No new kills on page %d for character %d, stopping", page, characterID)
			break
		}

		page++
		time.Sleep(1 * time.Second) // Add a delay to avoid hitting rate limits
	}

	log.Printf("Finished fetching kills for character %d. Total new kills: %d", characterID, totalNewKills)
}

func FetchAllKillsForCharacter(characterID int64) {
	log.Printf("Starting full kill fetch for character %d", characterID)
	fetchKillsForCharacter(characterID)
	log.Printf("Finished full kill fetch for character %d", characterID)
}

func FetchKillsForCharacter(characterID int64) {
	lastKillTime, err := db.GetLastKillTimeForCharacter(characterID)
	if err != nil {
		log.Printf("Error getting last kill time for character %d: %v", characterID, err)
		return
	}
	log.Printf("Last kill time for character %d: %v", characterID, lastKillTime)

	page := 1
	totalNewKills := 0

	for {
		log.Printf("Fetching page %d for character %d", page, characterID)
		url := fmt.Sprintf("https://zkillboard.com/api/characterID/%d/page/%d/", characterID, page)

		var zkillResponse []struct {
			KillmailID   int64  `json:"killmail_id"`
			KillmailTime string `json:"killmail_time"`
			ZKB          struct {
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

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching kills for character %d: %v", characterID, err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}

		err = json.Unmarshal(body, &zkillResponse)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			return
		}

		newKills := 0
		for _, zkill := range zkillResponse {
			killTime, err := time.Parse("2006-01-02T15:04:05Z", zkill.KillmailTime)
			if err != nil {
				log.Printf("Error parsing kill time for kill %d: %v", zkill.KillmailID, err)
				continue
			}

			if killTime.After(lastKillTime) {
				kill := models.Kill{
					KillmailID:     zkill.KillmailID,
					CharacterID:    characterID,
					KillTime:       killTime,
					LocationID:     zkill.ZKB.LocationID,
					Hash:           zkill.ZKB.Hash,
					FittedValue:    zkill.ZKB.FittedValue,
					DroppedValue:   zkill.ZKB.DroppedValue,
					DestroyedValue: zkill.ZKB.DestroyedValue,
					TotalValue:     zkill.ZKB.TotalValue,
					Points:         zkill.ZKB.Points,
					NPC:            zkill.ZKB.NPC,
					Solo:           zkill.ZKB.Solo,
					Awox:           zkill.ZKB.Awox,
				}

				// Fetch additional data from ESI
				esiKill, err := services.FetchKillmailFromESI(kill.KillmailID, kill.Hash)
				if err != nil {
					log.Printf("Error fetching killmail %d from ESI: %v", kill.KillmailID, err)
					continue
				}

				// Combine zKillboard and ESI data
				kill.SolarSystemID = esiKill.SolarSystemID
				kill.Victim = esiKill.Victim
				kill.Attackers = esiKill.Attackers

				// Insert the kill
				err = db.InsertKill(&kill)
				if err != nil {
					log.Printf("Error inserting kill %d: %v", kill.KillmailID, err)
				} else {
					newKills++
					totalNewKills++
				}
			}
		}

		log.Printf("Inserted %d new kills for character %d on page %d", newKills, characterID, page)

		if newKills == 0 {
			log.Printf("No new kills on page %d for character %d, stopping", page, characterID)
			break
		}

		page++
		time.Sleep(1 * time.Second) // Add a delay to avoid hitting rate limits
	}

	log.Printf("Finished fetching kills for character %d. Total new kills: %d", characterID, totalNewKills)
}
