package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"slantbit/internal/repository"
)

type createOrderReq struct {
	CourseID int64  `json:"courseId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Name == "" || !strings.Contains(req.Email, "@") {
		writeError(w, http.StatusBadRequest, "name and valid email required")
		return
	}

	course, err := h.store.GetCourse(r.Context(), req.CourseID)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "course not found")
		return
	}
	if err != nil {
		log.Println("get course:", err)
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	rzpOrder, err := h.rzp.Order.Create(map[string]interface{}{
		"amount":   course.PricePaise,
		"currency": "INR",
		"receipt":  fmt.Sprintf("rcpt_%d", time.Now().UnixNano()),
	}, nil)
	if err != nil {
		log.Println("razorpay create order:", err)
		writeError(w, http.StatusBadGateway, "payment provider error")
		return
	}
	rzpID, ok := rzpOrder["id"].(string)
	if !ok {
		writeError(w, http.StatusBadGateway, "payment provider error")
		return
	}

	if err := h.store.CreateOrder(r.Context(), course.ID, req.Name, req.Email, rzpID, course.PricePaise); err != nil {
		log.Println("save order:", err)
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"orderId": rzpID,
		"amount":  course.PricePaise,
		"key":     h.cfg.RazorpayKeyID,
	})
}