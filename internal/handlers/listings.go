package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Siddharth9101/olx-api/internal/middlewares"
)

type listing struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price string `json:"price"`
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
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	listings := []listing{}
	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("listing rows scan error", "request_id", requestId, "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		lh.logger.Error("listing rows error", "request_id", requestId, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
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
		lh.logger.Error("listings delete error", "request_id", requestId, "listing_id", id, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	lh.logger.Info("listing deleted successfully", "request_id", requestId, "listing_id", id)

	w.WriteHeader(http.StatusNoContent)
}
