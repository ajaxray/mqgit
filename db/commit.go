package db

import "encoding/json"

type Commit struct {
	ID         int64  `json:"ID"`
	Message    string `json:"Message,omitempty"`
	GitHash    string `json:"GitHash,omitempty"`
	GitMessage string `json:"GitMessage,omitempty"`
	Sql        string `json:"SQL"`
}

func (c Commit) ToJSON() (string, error) {
	if data, err := json.Marshal(c); err == nil {
		return string(data), nil
	} else {
		return "Error in JSON Parsing", err
	}
}

func (c *Commit) FromJSON(data string) error {
	err := json.Unmarshal([]byte(data), c)
	return err
}
