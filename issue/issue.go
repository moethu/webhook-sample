package issue

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Issue sample payload for webhooks
type Issue struct {
	Priority   int       `json:"priority"`
	Code       int       `json:"code"`
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Time       time.Time `json:"timestamp"`
	System     uuid.UUID `json:"system"`
	Confidence float32   `json:"confidence"`
}

// NewIssue instanciates a new default issue
func NewIssue() Issue {
	id := uuid.NewV4()
	i := Issue{
		ID:         id,
		System:     uuid.Nil,
		Time:       time.Now(),
		Content:    "None",
		Title:      "Untitled",
		Confidence: 0.0,
		Priority:   0,
		Code:       0}
	return i
}
