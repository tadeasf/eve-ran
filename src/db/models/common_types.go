package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type IntArray []int

func (a *IntArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var arr []int
	err := json.Unmarshal(bytes, &arr)
	*a = IntArray(arr)
	return err
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (p *Position) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, &p)
	return err
}
