package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}
type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

//Constructor pattern

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

// Closure Factory - technical term
// this wrapping a function inside another function is the use case of closure
// and this is dependency injection in go

// func Listings(db *sql.DB) http.HandlerFunc {

func (lh ListingHandler) Listings(w http.ResponseWriter, r *http.Request) {
	//request scoped context
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx,
		`SELECT id,title,description,price,city,created_at 
			FROM listings 
			ORDER BY created_at DESC 
			LIMIT 100`)
	if err != nil {
		lh.logger.Error("db.QueryContext", "err", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close() //closes the connection established by the rows

	listings := []listing{}

	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("rows.scan:", "err", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		lh.logger.Error("rows.err: ", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listings)
}

// }

func (lh ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	_, err := lh.db.ExecContext(ctx, `DELETE FROM listing WHERE id = $1`, id)
	if err != nil {
		// log.Printf("delete db.Exec Fail : %v", err)
		lh.logger.Error("delete failed", "listing_id", id, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
