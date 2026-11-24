package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/karan-singh-bisht/Dealzo-api/internal/httpx"
	"github.com/karan-singh-bisht/Dealzo-api/internal/middleware"
	"github.com/karan-singh-bisht/Dealzo-api/internal/storage"
)

type UploadHandler struct {
	logger *slog.Logger
	store  *storage.Client
}

func NewUploadHandler(logger *slog.Logger, store *storage.Client) *UploadHandler {
	return &UploadHandler{
		logger: logger,
		store:  store,
	}
}

func (uh UploadHandler) Presign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIdFromContext(ctx) //returned value from middleware
	userID, err := middleware.UserIdFromContext(ctx)
	log := uh.logger.With("request_id", requestId)

	if err != nil {
		uh.logger.Error("No user id found in context", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	var req PresignRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Decoding failed", "err", err)
		httpx.Error(w, http.StatusBadRequest, "something went wrong", httpx.CodeValidationFailed)
		return
	}

	if len(req.Files) == 0 {
		log.Error("Files are empty", "err", err)
		httpx.Error(w, http.StatusBadRequest, "Files must not be empty", httpx.CodeValidationFailed)
		return
	}

	if len(req.Files) > maxImagePerListing {
		log.Error(fmt.Sprintf("Number of images exceed %d", maxImagePerListing), "err", err)
		httpx.Error(w, http.StatusBadRequest, fmt.Sprintf("Images cannot exceed %d", maxImagePerListing), httpx.CodeValidationFailed)
		return
	}

	for _, f := range req.Files {
		_, ok := allowedContentTypes[f.ContentType]
		if !ok {
			log.Error(fmt.Sprintf("%s extension is not allowed", f.ContentType), "err", err)
			httpx.Error(w, http.StatusBadRequest, "only image/jpeg, image/png, image/webp are allowed", httpx.CodeMalformedJSON)
			return
		}

		if f.SizeBytes <= 0 || f.SizeBytes > maxImageBytes {
			log.Error("file exceeds max byte size", "err", err)
			httpx.Error(w, http.StatusBadRequest, fmt.Sprintf("File exceeds max size of %d bytes", maxImageBytes), httpx.CodeMalformedJSON)
			return
		}
		//Actually Presigning
		uploadKey := mintUploadKey(uuid.UUID(userID), allowedContentTypes[f.ContentType])
		fmt.Println(uploadKey)

	}

	w.Write([]byte("Hello"))
}
