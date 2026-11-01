package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/karan-singh-bisht/Dealzo-api/internal/config"
	"github.com/karan-singh-bisht/Dealzo-api/internal/httpx"
	"github.com/karan-singh-bisht/Dealzo-api/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

var dummyHash = []byte("$2a$10$KIhhQ9B9yG5Y9mNlNZYPzuJfkb0Y08NAM7j1ABPv2HAczFISuptjK")

type user struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_at time.Time `json:"created_at"`
}

type AuthHandler struct {
	db     *sql.DB
	logger *slog.Logger
	cfg    config.Config
}

func NewAuthHandler(db *sql.DB, logger *slog.Logger, cfg config.Config) *AuthHandler {
	return &AuthHandler{
		db:     db,
		logger: logger,
		cfg:    cfg,
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
		httpx.ValidationError(w, http.StatusBadRequest, verr.Error(), httpx.CodeValidationFailed, verr.Field)
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

func (ah AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	request_id := middleware.RequestIdFromContext(ctx)
	ah.logger.With("request_id", request_id)

	var req SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.logger.Error("Decoding Failed at user sign in", "err", err)
		httpx.Error(w, http.StatusBadRequest, "Bad Request", httpx.CodeValidationFailed)
		return
	}

	if err := req.Validate(); err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		httpx.ValidationError(w, http.StatusBadRequest, verr.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	row := ah.db.QueryRowContext(ctx,
		`SELECT id,email,password 
		FROM users 
		WHERE email = $1`, req.Email)

	//2 error methods (As & Is)

	var u user
	if err := row.Scan(&u.ID, &u.Email, &u.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//We are doing this to add a layer of security because hackers can judge whether their password is wrong or email
			//is wrong based on response time so we are trying to keep response time the same
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
			httpx.Error(w, http.StatusUnauthorized, "Wrong Email or Password", httpx.CodeUnauthenticated)
			return
		}
		ah.logger.Error("Failed to scan", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		ah.logger.Warn("Wrong Credentials", "err", err, "user_id", u.ID)
		httpx.Error(w, http.StatusUnauthorized, "Wrong Email or Password", httpx.CodeUnauthenticated)
		return
	}

	tokenTTL := 24 * time.Hour

	claims := jwt.RegisteredClaims{
		Subject:   u.ID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(ah.cfg.JWT_token))

	if err != nil {
		ah.logger.Error("Token Signing failed", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	out := SigninResponse{
		Token:     signedToken,
		ExpiresIn: int(tokenTTL.Seconds()),
	}

	ah.logger.Info("User Logged In", "user_id", u.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(out)
}
