package main

import (
	"encoding/json"
	"fmt"
)

type Item struct {
	Question        []byte
	Value           [][]byte
	Content         [][]byte
	ValueStartIndex int
	Index           int
}

func (item Item) String() {
	fmt.Printf("Count: %5d, Question: %s, Value: %s \n", item.Index+1, item.Question, item.Value)
}

func (item Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Question string   `json:"question"`
		Value    []string `json:"value"`
		Content  []string `json:"content"`
	}{
		Question: string(item.Question),
		Value:    convertByteArrayToString(item.Value),
		Content:  convertByteArrayToString(item.Content),
	})
}

func convertByteArrayToString(b [][]byte) []string {
	strs := make([]string, len(b))
	for i, v := range b {
		strs[i] = string(v)
	}
	return strs
}
