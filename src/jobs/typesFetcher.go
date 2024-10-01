package jobs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/services"
)

const (
	baseURL = "https://esi.evetech.net/latest"
)

func FetchAndUpdateTypes() {
	fetchAndUpdateRegions()
	fetchAndUpdateConstellations()
	fetchAndUpdateSystems()
	fetchAndUpdateItems()
}

func fetchAndUpdateRegions() {
	regions, err := services.FetchAllRegions(10)
	if err != nil {
		fmt.Printf("Error fetching regions: %v\n", err)
		return
	}

	for _, region := range regions {
		err := db.UpsertRegion(region)
		if err != nil {
			fmt.Printf("Error upserting region %d: %v\n", region.RegionID, err)
		}
	}
}

func fetchAndUpdateConstellations() {
	url := fmt.Sprintf("%s/universe/constellations/", baseURL)
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
}

func fetchAndSaveConstellation(id int) {
	url := fmt.Sprintf("%s/universe/constellations/%d/", baseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching constellation %d: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var constellation models.Constellation
	json.Unmarshal(body, &constellation)

	db.UpsertConstellation(&constellation)
}

func fetchAndUpdateSystems() {
	url := fmt.Sprintf("%s/universe/systems/", baseURL)
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
}

func fetchAndSaveSystem(id int) {
	url := fmt.Sprintf("%s/universe/systems/%d/", baseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching system %d: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var system models.System
	json.Unmarshal(body, &system)

	db.UpsertSystem(&system)
}

func fetchAndUpdateItems() {
	url := fmt.Sprintf("%s/universe/types/", baseURL)
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
}

func fetchAndSaveItem(id int) {
	url := fmt.Sprintf("%s/universe/types/%d/", baseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching item %d: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var item models.ESIItem
	json.Unmarshal(body, &item)

	db.UpsertESIItem(&item)
}

func fetchIDs(url string) []int {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching IDs: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ids []int
	json.Unmarshal(body, &ids)

	return ids
}
