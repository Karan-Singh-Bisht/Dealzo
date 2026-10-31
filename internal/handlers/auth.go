package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/karan-singh-bisht/Dealzo-api/internal/httpx"
	"github.com/karan-singh-bisht/Dealzo-api/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewAuthHandler(db *sql.DB, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		db:     db,
		logger: logger,
	}
}

func (ah AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIdFromContext(ctx)

	log := ah.logger.With("request_id", requestId) //ingests request id with each log and return a new logger

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to create user", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "Something went wrong", httpx.CodeValidationFailed)
		return
	}

	if err := req.AuthValidate(); err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		httpx.ValidationError(w, http.StatusBadRequest, err.Error(), httpx.CodeValidationFailed, verr.Field)
	}

	//cost -> how slow you want your hashing function - 2 ^ cost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		log.Error("hashing failed", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	row := ah.db.QueryRowContext(ctx, `
	INSERT INTO users
	(name,email,password) 
	VALUES ($1,$2,$3) RETURNING id,created_at`, req.Name, req.Email, hashedPassword)

	var out CreateUserResponse

	if err := row.Scan(&out.ID, &out.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		//status conflict - 409 when trying to store duplicate data
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			httpx.Error(w, http.StatusConflict, "Email Already Taken", httpx.CodeConflict)
			return
		}
		log.Error("Failed to scan", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	log.Info("New User Created", "user_id", out.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(out)
}
