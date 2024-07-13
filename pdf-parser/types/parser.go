package types

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type InfoPDF struct {
	URL  *url.URL
	Name string
}

type Enrollment struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	URL       string    `json:"url,omitempty"`
	Content   []byte    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
