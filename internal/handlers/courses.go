package handlers

import (
	"log"
	"net/http"
)

func (h *Handler) ListCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.store.ListCourses(r.Context())
	if err != nil {
		log.Println("list courses:", err)
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	writeJSON(w, http.StatusOK, courses)
}