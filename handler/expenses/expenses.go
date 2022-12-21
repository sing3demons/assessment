package expenses

import (
	"database/sql"

	"github.com/lib/pq"
)

type handler struct {
	DB *sql.DB
}

type Err struct {
	Message string `json:"message"`
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}

type NewsExpenses struct {
	ID     int            `json:"id"`
	Title  string         `json:"title"`
	Amount float64        `json:"amount"`
	Note   string         `json:"note"`
	Tags   pq.StringArray `json:"tags"`
}
