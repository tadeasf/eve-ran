package jobs

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/services"
)

const (
	baseURL = "https://esi.evetech.net/latest"
)

func FetchAndUpdateTypes() {
	log.Println("Starting FetchAndUpdateTypes job")
	fetchAndUpdateRegions()
	fetchAndUpdateConstellations()
	fetchAndUpdateSystems()
	fetchAndUpdateItems()
	log.Println("Finished FetchAndUpdateTypes job")
}

func fetchAndUpdateRegions() {
	log.Println("Fetching and updating regions")
	regions, err := services.FetchAllRegions(10)
	if err != nil {
		log.Printf("Error fetching regions: %v", err)
		return
	}

	for _, region := range regions {
		err := db.UpsertRegion(region)
		if err != nil {
			log.Printf("Error upserting region %d: %v", region.RegionID, err)
		}
	}
	log.Println("Finished fetching and updating regions")
}

func fetchAndUpdateConstellations() {
	log.Println("Fetching and updating constellations")
	url := baseURL + "/universe/constellations/"
	ids := fetchIDs(url)

	existingConstellations, _ := db.GetAllConstellations()
	existingMap := make(map[int]bool)
	for _, constellation := range existingConstellations {
		existingMap[constellation.ConstellationID] = true
	}

	for _, id := range ids {
		if !existingMap[id] {
			fetchAndSaveConstellation(id)
		}
	}
	log.Println("Finished fetching and updating constellations")
}

func fetchAndSaveConstellation(id int) {
	url := baseURL + "/universe/constellations/" + strconv.Itoa(id) + "/"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching constellation %d: %v", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var constellation models.Constellation
	json.Unmarshal(body, &constellation)

	err = db.UpsertConstellation(&constellation)
	if err != nil {
		log.Printf("Error upserting constellation %d: %v", id, err)
	}
}

func fetchAndUpdateSystems() {
	log.Println("Fetching and updating systems")
	url := baseURL + "/universe/systems/"
	ids := fetchIDs(url)

	existingSystems, _ := db.GetAllSystems()
	existingMap := make(map[int]bool)
	for _, system := range existingSystems {
		existingMap[system.SystemID] = true
	}

	for _, id := range ids {
		if !existingMap[id] {
			fetchAndSaveSystem(id)
		}
	}
	log.Println("Finished fetching and updating systems")
}

func fetchAndSaveSystem(id int) {
	url := baseURL + "/universe/systems/" + strconv.Itoa(id) + "/"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching system %d: %v", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var system models.System
	json.Unmarshal(body, &system)

	err = db.UpsertSystem(&system)
	if err != nil {
		log.Printf("Error upserting system %d: %v", id, err)
	}
}

func fetchAndUpdateItems() {
	log.Println("Fetching and updating items")
	url := baseURL + "/universe/types/"
	ids := fetchIDs(url)

	existingItems, _ := db.GetAllESIItems()
	existingMap := make(map[int]bool)
	for _, item := range existingItems {
		existingMap[item.TypeID] = true
	}

	for _, id := range ids {
		if !existingMap[id] {
			fetchAndSaveItem(id)
		}
	}
	log.Println("Finished fetching and updating items")
}

func fetchAndSaveItem(id int) {
	if id == 0 {
		log.Printf("Skipping item with ID 0")
		return
	}
	url := baseURL + "/universe/types/" + strconv.Itoa(id) + "/"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching item %d: %v", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var item models.ESIItem
	json.Unmarshal(body, &item)

	err = db.UpsertESIItem(&item)
	if err != nil {
		log.Printf("Error upserting item %d: %v", id, err)
	}
}

func fetchIDs(url string) []int {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching IDs from %s: %v", url, err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ids []int
	json.Unmarshal(body, &ids)

	return ids
}
