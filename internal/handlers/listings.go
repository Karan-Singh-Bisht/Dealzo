package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/karan-singh-bisht/Dealzo-api/internal/httpx"
	"github.com/karan-singh-bisht/Dealzo-api/internal/middleware"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
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
	requestId := middleware.RequestIdFromContext(ctx) //returned value from middleware
	id := r.PathValue("id")

	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		// log.Printf("delete db.Exec Fail : %v", err)
		lh.logger.Error("delete failed", "listing_id", id, "requestId", requestId, "err", err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (lh ListingHandler) CreateListing(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestId := middleware.RequestIdFromContext(ctx) //returned value from middleware

	var req listing
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("Decoding failed", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "Failed to decode", httpx.CodeInvalidID)
		return
	}

	row := lh.db.QueryRowContext(ctx,
		`INSERT INTO listings 
		(title,description,price,city)
		VALUES ($1,$2,$3,$4) RETURNING id
	`,
		req.Title, req.Description, req.Price, req.City,
	)

	var id string

	if err := row.Scan(&id); err != nil {
		lh.logger.Error("Row Scan Failed", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("listing created", "request_id", requestId, "listing_id", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}
