package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Siddharth9101/olx-api/internal/httpx"
	"github.com/Siddharth9101/olx-api/internal/middlewares"
)

type listing struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price int64 `json:"price"`
	City string `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

type ListingHandler struct {
	db *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db: db,
		logger: logger,
	}
}

func (lh ListingHandler) List(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	requestId := middlewares.RequestIdFromCtx(ctx)
	rows, err := lh.db.QueryContext(ctx,
		`SELECT id, title, description, price, city, created_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT 100`)
	if err != nil {
		lh.logger.Error("listing query error", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	defer rows.Close()
	
	listings := []listing{}
	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("listing rows scan error", "request_id", requestId, "err", err)
			httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
			return
		}
		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		lh.logger.Error("listing rows error", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("listings fetched", "request_id", requestId, "total", len(listings))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listings)
}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()
	requestId := middlewares.RequestIdFromCtx(ctx)

	_, err := lh.db.ExecContext(ctx, `DELETE FROM listing WHERE id = $1`, id)
	if err != nil {
		lh.logger.Error("listing delete error", "request_id", requestId, "listing_id", id, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	} 

	lh.logger.Info("listing deleted successfully", "request_id", requestId, "listing_id", id)

	w.WriteHeader(http.StatusNoContent)
}

func (lh ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIdFromCtx(ctx)
	var req listing
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid req body", httpx.CodeMalformedJson)
		return
	}

	row := lh.db.QueryRowContext(ctx, `
	INSERT INTO listings (title, description, price, city) VALUES ($1, $2, $3, $4) RETURNING id`, req.Title, req.Description, req.Price, req.City)

	var id string
	if err := row.Scan(&id); err != nil{
		lh.logger.Error("failed to insert", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("listing created", "request_id", requestId, "listing_id", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}