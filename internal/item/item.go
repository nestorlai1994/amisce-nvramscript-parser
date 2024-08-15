package item

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

func (item *Item) UnmarshalJSON(data []byte) error {
	var v struct {
		Question string   `json:"question"`
		Value    []string `json:"value"`
		Content  []string `json:"content"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	// fmt.Println("v: ", v.Question)
	item.Question = []byte(v.Question)
	item.Value = convertStringToByteArray(v.Value)
	item.Content = convertStringToByteArray(v.Content)
	return nil
}

func convertStringToByteArray(s []string) [][]byte {
	b := make([][]byte, len(s))
	for i, v := range s {
		b[i] = []byte(v)
	}
	return b
}

func (item Item) IsValueEqual(value [][]byte) bool {
	if len(item.Value) != len(value) {
		return false
	}
	for i, item := range item.Value {
		if string(item) != string(value[i]) {
			return false
		}
	}
	return true
}
