package jobs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

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
	baseURL := baseURL + "/universe/types/"

	existingItems, _ := db.GetAllESIItems()
	existingMap := make(map[int]bool)
	for _, item := range existingItems {
		existingMap[item.TypeID] = true
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20) // Limit to 20 concurrent requests
	itemIDsChan := make(chan int, 100)

	// Start worker goroutines
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range itemIDsChan {
				semaphore <- struct{}{}
				fetchAndSaveItem(id)
				<-semaphore
			}
		}()
	}

	page := 1
	for {
		ids, err := fetchItemIDsWithPagination(baseURL, page)
		if err != nil {
			if err.Error() == "requested page does not exist" {
				log.Println("Reached the end of item pages")
				break
			}
			log.Printf("Error fetching item IDs for page %d: %v", page, err)
			break
		}

		for _, id := range ids {
			if !existingMap[id] {
				itemIDsChan <- id
			}
		}

		page++
		time.Sleep(100 * time.Millisecond) // Small delay to avoid hitting rate limits
	}

	close(itemIDsChan)
	wg.Wait()

	log.Println("Finished fetching and updating items")
}

func fetchItemIDsWithPagination(baseURL string, page int) ([]int, error) {
	url := fmt.Sprintf("%s?datasource=tranquility&page=%d", baseURL, page)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("requested page does not exist")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ids []int
	err = json.Unmarshal(body, &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func fetchAndSaveItem(id int) {
	if id == 0 {
		log.Printf("Skipping item with ID 0")
		return
	}
	url := fmt.Sprintf("%s/universe/types/%d/?datasource=tranquility&language=en", baseURL, id)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for item %d: %v", id, err)
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching item %d: %v", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var item models.ESIItem
	err = json.Unmarshal(body, &item)
	if err != nil {
		log.Printf("Error unmarshaling item %d: %v", id, err)
		return
	}

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
