package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/tadeasf/eve-ran/src/db/models"
)

var DB *sql.DB

func InitDB() {
	connStr := "host=localhost port=5435 user=eve password=eve dbname=eve sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")

	// Create tables if they don't exist
	createTables()
}

func createTables() {
	tables := []interface{}{
		models.Character{},
		models.Kill{},
	}

	for _, table := range tables {
		createTableFromModel(table)
	}
}

func createTableFromModel(model interface{}) {
	t := reflect.TypeOf(model)
	tableName := strings.ToLower(t.Name()) + "s"

	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			parts := strings.Split(dbTag, ",")
			if len(parts) < 2 {
				log.Printf("Warning: Invalid db tag for field %s, skipping", field.Name)
				continue
			}
			columnName := parts[0]
			columnType := parts[1]
			constraints := ""
			if len(parts) > 2 {
				constraints = strings.Join(parts[2:], " ")
			}
			columns = append(columns, fmt.Sprintf("%s %s %s", columnName, columnType, constraints))
		}
	}

	if len(columns) == 0 {
		log.Printf("Warning: No valid columns found for table %s, skipping table creation", tableName)
		return
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n\t%s\n);", tableName, strings.Join(columns, ",\n\t"))

	_, err := DB.Exec(query)
	if err != nil {
		log.Printf("Error creating table %s: %v", tableName, err)
	} else {
		log.Printf("Table %s created or already exists", tableName)
	}
}

func InsertCharacter(character *models.Character) error {
	_, err := DB.Exec("INSERT INTO characters (id) VALUES ($1) ON CONFLICT (id) DO NOTHING", character.ID)
	return err
}

func InsertKill(kill *models.Kill) error {
	victimJSON, err := json.Marshal(kill.Victim)
	if err != nil {
		return err
	}
	attackersJSON, err := json.Marshal(kill.Attackers)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		INSERT INTO kills (
			killmail_id, character_id, killmail_time, solar_system_id,
			location_id, hash, fitted_value, dropped_value, destroyed_value,
			total_value, points, npc, solo, awox, victim, attackers
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (killmail_id) DO UPDATE SET
			character_id = EXCLUDED.character_id,
			killmail_time = EXCLUDED.killmail_time,
			solar_system_id = EXCLUDED.solar_system_id,
			location_id = EXCLUDED.location_id,
			hash = EXCLUDED.hash,
			fitted_value = EXCLUDED.fitted_value,
			dropped_value = EXCLUDED.dropped_value,
			destroyed_value = EXCLUDED.destroyed_value,
			total_value = EXCLUDED.total_value,
			points = EXCLUDED.points,
			npc = EXCLUDED.npc,
			solo = EXCLUDED.solo,
			awox = EXCLUDED.awox,
			victim = EXCLUDED.victim,
			attackers = EXCLUDED.attackers`,
		kill.KillmailID, kill.CharacterID, kill.KillTime, kill.SolarSystemID,
		kill.LocationID, kill.Hash, kill.FittedValue, kill.DroppedValue,
		kill.DestroyedValue, kill.TotalValue, kill.Points,
		kill.NPC, kill.Solo, kill.Awox, victimJSON, attackersJSON)
	return err
}

func GetKillsForCharacter(characterID int64, page, pageSize int) ([]models.Kill, error) {
	offset := (page - 1) * pageSize
	rows, err := DB.Query(`
		SELECT killmail_id, character_id, killmail_time, solar_system_id,
			   location_id, hash, fitted_value, dropped_value, destroyed_value,
			   total_value, points, npc, solo, awox, victim, attackers
		FROM kills
		WHERE character_id = $1
		ORDER BY killmail_time DESC
		LIMIT $2 OFFSET $3
	`, characterID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kills []models.Kill
	for rows.Next() {
		var kill models.Kill
		var victimJSON, attackersJSON []byte
		err := rows.Scan(
			&kill.KillmailID, &kill.CharacterID, &kill.KillTime, &kill.SolarSystemID,
			&kill.LocationID, &kill.Hash, &kill.FittedValue, &kill.DroppedValue,
			&kill.DestroyedValue, &kill.TotalValue, &kill.Points,
			&kill.NPC, &kill.Solo, &kill.Awox, &victimJSON, &attackersJSON,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(victimJSON, &kill.Victim)
		json.Unmarshal(attackersJSON, &kill.Attackers)
		kills = append(kills, kill)
	}

	return kills, nil
}

func GetTotalKillsForCharacter(characterID int64) (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM kills WHERE character_id = $1", characterID).Scan(&count)
	return count, err
}

func GetAllCharacters() ([]models.Character, error) {
	log.Println("Fetching all characters from database")
	rows, err := DB.Query("SELECT id FROM characters")
	if err != nil {
		log.Printf("Error querying characters: %v", err)
		return nil, err
	}
	defer rows.Close()

	var characters []models.Character
	for rows.Next() {
		var char models.Character
		if err := rows.Scan(&char.ID); err != nil {
			log.Printf("Error scanning character: %v", err)
			return nil, err
		}
		characters = append(characters, char)
	}
	log.Printf("Found %d characters in database", len(characters))
	return characters, nil
}

func GetLastKillTimeForCharacter(characterID int64) (time.Time, error) {
	var lastKillTime sql.NullTime
	err := DB.QueryRow("SELECT MAX(killmail_time) FROM kills WHERE character_id = $1", characterID).Scan(&lastKillTime)
	if err != nil {
		return time.Time{}, err
	}
	if !lastKillTime.Valid {
		return time.Time{}, nil
	}
	return lastKillTime.Time, nil
}

func GetKillByKillmailID(killmailID int64) (*models.Kill, error) {
	var kill models.Kill
	var victimJSON, attackersJSON []byte

	err := DB.QueryRow(`
		SELECT killmail_id, character_id, killmail_time, solar_system_id,
			   location_id, hash, fitted_value, dropped_value, destroyed_value,
			   total_value, points, npc, solo, awox, victim, attackers
		FROM kills
		WHERE killmail_id = $1`, killmailID).Scan(
		&kill.KillmailID, &kill.CharacterID, &kill.KillTime, &kill.SolarSystemID,
		&kill.LocationID, &kill.Hash, &kill.FittedValue, &kill.DroppedValue,
		&kill.DestroyedValue, &kill.TotalValue, &kill.Points,
		&kill.NPC, &kill.Solo, &kill.Awox, &victimJSON, &attackersJSON)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(victimJSON, &kill.Victim)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(attackersJSON, &kill.Attackers)
	if err != nil {
		return nil, err
	}

	return &kill, nil
}
