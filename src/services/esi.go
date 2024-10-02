package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
)

const esiBaseURL = "https://esi.evetech.net/latest"

func FetchRegionIDs() ([]int, error) {
	url := fmt.Sprintf("%s/universe/regions/?datasource=tranquility", esiBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var regionIDs []int
	err = json.Unmarshal(body, &regionIDs)
	return regionIDs, err
}

func FetchRegionInfo(regionID int) (*models.Region, error) {
	url := fmt.Sprintf("%s/universe/regions/%d/?datasource=tranquility&language=en", esiBaseURL, regionID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var region models.Region
	err = json.Unmarshal(body, &region)
	if err != nil {
		return nil, err
	}

	// Ensure Constellations is initialized as an empty slice if it's null
	if region.Constellations == nil {
		region.Constellations = []int{}
	}

	return &region, nil
}

func FetchSystemIDs() ([]int, error) {
	url := fmt.Sprintf("%s/universe/systems/?datasource=tranquility", esiBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var systemIDs []int
	err = json.Unmarshal(body, &systemIDs)
	return systemIDs, err
}

func FetchSystemInfo(systemID int) (*models.System, error) {
	url := fmt.Sprintf("%s/universe/systems/%d/?datasource=tranquility&language=en", esiBaseURL, systemID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var system models.System
	err = json.Unmarshal(body, &system)
	return &system, err
}

func FetchConstellationIDs() ([]int, error) {
	url := fmt.Sprintf("%s/universe/constellations/?datasource=tranquility", esiBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var constellationIDs []int
	err = json.Unmarshal(body, &constellationIDs)
	return constellationIDs, err
}

func FetchConstellationInfo(constellationID int) (*models.Constellation, error) {
	url := fmt.Sprintf("%s/universe/constellations/%d/?datasource=tranquility&language=en", esiBaseURL, constellationID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var constellation models.Constellation
	err = json.Unmarshal(body, &constellation)
	return &constellation, err
}

func FetchItemIDs() ([]int, error) {
	var allItemIDs []int
	page := 1
	for {
		url := fmt.Sprintf("%s/universe/types/?datasource=tranquility&page=%d", esiBaseURL, page)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var itemIDs []int
		err = json.Unmarshal(body, &itemIDs)
		if err != nil {
			return nil, err
		}

		if len(itemIDs) == 0 {
			break
		}

		allItemIDs = append(allItemIDs, itemIDs...)
		page++
	}

	return allItemIDs, nil
}

func FetchItemInfo(itemID int) (*models.ESIItem, error) {
	url := fmt.Sprintf("%s/universe/types/%d/?datasource=tranquility&language=en", esiBaseURL, itemID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item models.ESIItem
	err = json.Unmarshal(body, &item)
	return &item, err
}

func FetchAllItems(concurrency int) ([]*models.ESIItem, error) {
	itemIDs, err := FetchItemIDs()
	if err != nil {
		return nil, err
	}

	items := make([]*models.ESIItem, 0, len(itemIDs))
	itemChan := make(chan *models.ESIItem, len(itemIDs))
	errChan := make(chan error, len(itemIDs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, itemID := range itemIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			item, err := FetchItemInfo(id)
			if err != nil {
				errChan <- err
				return
			}
			itemChan <- item
		}(itemID)
	}

	go func() {
		wg.Wait()
		close(itemChan)
		close(errChan)
	}()

	for item := range itemChan {
		items = append(items, item)
	}

	if len(errChan) > 0 {
		return items, <-errChan
	}

	return items, nil
}

func FetchAllRegions(concurrency int) ([]*models.Region, error) {
	regionIDs, err := FetchRegionIDs()
	if err != nil {
		return nil, err
	}

	regions := make([]*models.Region, 0, len(regionIDs))
	regionChan := make(chan *models.Region, len(regionIDs))
	errChan := make(chan error, len(regionIDs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, regionID := range regionIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			region, err := FetchRegionInfo(id)
			if err != nil {
				errChan <- err
				return
			}
			regionChan <- region
		}(regionID)
	}

	go func() {
		wg.Wait()
		close(regionChan)
		close(errChan)
	}()

	for region := range regionChan {
		regions = append(regions, region)
	}

	if len(errChan) > 0 {
		return regions, <-errChan
	}

	return regions, nil
}

func FetchAllConstellations(concurrency int) ([]*models.Constellation, error) {
	constellationIDs, err := FetchConstellationIDs()
	if err != nil {
		return nil, err
	}

	constellations := make([]*models.Constellation, 0, len(constellationIDs))
	constellationChan := make(chan *models.Constellation, len(constellationIDs))
	errChan := make(chan error, len(constellationIDs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, constellationID := range constellationIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			constellation, err := FetchConstellationInfo(id)
			if err != nil {
				errChan <- err
				return
			}
			constellationChan <- constellation
		}(constellationID)
	}

	go func() {
		wg.Wait()
		close(constellationChan)
		close(errChan)
	}()

	for constellation := range constellationChan {
		constellations = append(constellations, constellation)
	}

	if len(errChan) > 0 {
		return constellations, <-errChan
	}

	return constellations, nil
}

func FetchAllSystems(concurrency int) ([]*models.System, error) {
	systemIDs, err := FetchSystemIDs()
	if err != nil {
		return nil, err
	}

	systems := make([]*models.System, 0, len(systemIDs))
	systemChan := make(chan *models.System, len(systemIDs))
	errChan := make(chan error, len(systemIDs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, systemID := range systemIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			system, err := FetchSystemInfo(id)
			if err != nil {
				errChan <- err
				return
			}
			systemChan <- system
		}(systemID)
	}

	go func() {
		wg.Wait()
		close(systemChan)
		close(errChan)
	}()

	for system := range systemChan {
		systems = append(systems, system)
	}

	if len(errChan) > 0 {
		return systems, <-errChan
	}

	return systems, nil
}

func FetchKillmailFromESI(killmailID int64, hash string) (*models.Kill, error) {
	maxRetries := 3
	baseDelay := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		kill, err := fetchKillmailFromESIWithRetry(killmailID, hash)
		if err == nil {
			return kill, nil
		}

		// Check if the error is a timeout error
		var esiError struct {
			Error   string `json:"error"`
			Timeout int    `json:"timeout"`
		}
		if jsonErr := json.Unmarshal([]byte(err.Error()), &esiError); jsonErr == nil && esiError.Error == "Timeout contacting tranquility" {
			delay := time.Duration(esiError.Timeout) * time.Second
			log.Printf("ESI timeout for killmail_id %d. Retrying in %v (attempt %d/%d)", killmailID, delay, attempt+1, maxRetries)
			time.Sleep(delay)
		} else if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(attempt+1)
			log.Printf("Retrying fetch for killmail_id %d in %v (attempt %d/%d)", killmailID, delay, attempt+1, maxRetries)
			time.Sleep(delay)
		} else {
			return nil, fmt.Errorf("failed to fetch killmail after %d attempts: %v", maxRetries, err)
		}
	}

	return nil, fmt.Errorf("unexpected error: should not reach this point")
}

func fetchKillmailFromESIWithRetry(killmailID int64, hash string) (*models.Kill, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/killmails/%d/%s/?datasource=tranquility", killmailID, hash)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "EVE Ran Application - GitHub: tadeasf/eve-ran")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading ESI response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", string(body))
	}

	var kill models.Kill
	err = json.Unmarshal(body, &kill)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling killmail: %v", err)
	}

	return &kill, nil
}
