package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
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
	db *sql.DB
}

//Constructor pattern

func NewListingHandler(db *sql.DB) *ListingHandler {
	return &ListingHandler{
		db: db,
	}
}

// func Listings(db *sql.DB) http.HandlerFunc {

// Closure Factory - technical term
// this wrapping a function inside another function is the use case of closure
// and this is dependency injection in go
func (lh ListingHandler) Listings(w http.ResponseWriter, r *http.Request) {
	//request scoped context
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx,
		`SELECT id,title,description,price,city,created_at 
			FROM listings 
			ORDER BY created_at DESC 
			LIMIT 100`)
	if err != nil {
		log.Printf("query: %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close() //closes the connection established by the rows

	listings := []listing{}

	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			log.Printf("rows.scan : %v", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows.err: %v", err)
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
	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		log.Printf("delete db.Exec Fail : %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
