package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"slantbit/internal/repository"
)

func (h *Handler) GetAccess(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	a, err := h.store.GetAccessByToken(r.Context(), token)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "invalid access link")
		return
	}
	if err != nil {
		log.Println("get access:", err)
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	pdfURL, err := h.storage.PresignGet(r.Context(), a.PdfKey, 10*time.Minute)
	if err != nil {
		log.Println("presign:", err)
		writeError(w, http.StatusInternalServerError, "could not generate download link")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"courseTitle":  a.CourseTitle,
		"pdfUrl":       pdfURL,
		"whatsappLink": a.WhatsappLink,
	})
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	var req struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Rating < 1 || req.Rating > 5 {
		writeError(w, http.StatusBadRequest, "rating must be between 1 and 5")
		return
	}
	req.Comment = strings.TrimSpace(req.Comment)
	if len(req.Comment) > 1000 {
		writeError(w, http.StatusBadRequest, "comment too long")
		return
	}

	a, err := h.store.GetAccessByToken(r.Context(), token)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "invalid access link")
		return
	}
	if err != nil {
		log.Println("get access:", err)
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	err = h.store.CreateReview(r.Context(), a.EnrollmentID, a.CourseID, req.Rating, req.Comment)
	if errors.Is(err, repository.ErrAlreadyReviewed) {
		writeError(w, http.StatusConflict, "you already reviewed this course")
		return
	}
	if err != nil {
		log.Println("create review:", err)
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}