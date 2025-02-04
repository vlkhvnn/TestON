package models

import "encoding/json"

type RecentChangeEvent struct {
	ID         json.Number `json:"id"`
	Type       string      `json:"type"`
	Title      string      `json:"title"`
	User       string      `json:"user"`
	Bot        bool        `json:"bot"`
	Minor      bool        `json:"minor"`
	Comment    string      `json:"comment"`
	Timestamp  int64       `json:"timestamp"`
	Wiki       string      `json:"wiki"`
	ServerName string      `json:"server_name"`
}
