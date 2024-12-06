package order

import "time"

type Order struct {
	ID         int64      `json:"ID"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	Type       string     `json:"type"`
	ModelName  string     `json:"modelName"`
	Desc       string     `json:"desc"`
	Status     Status     `json:"status"`
	Comment    string     `json:"comment,omitempty"`
	ClientData ClientData `json:"clientData"`
	Master     string     `json:"master,omitempty"`
}
