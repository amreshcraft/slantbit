package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (h *Handler) RazorpayWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}

	// raw body pe signature verify, JSON parse se pehle
	mac := hmac.New(sha256.New, []byte(h.cfg.RazorpayWebhookSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(r.Header.Get("X-Razorpay-Signature"))) {
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	var evt struct {
		Event   string `json:"event"`
		Payload struct {
			Payment struct {
				Entity struct {
					ID      string `json:"id"`
					OrderID string `json:"order_id"`
				} `json:"entity"`
			} `json:"payment"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(body, &evt); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if evt.Event != "payment.captured" {
		w.WriteHeader(http.StatusOK)
		return
	}

	e := evt.Payload.Payment.Entity
	res, err := h.store.CompletePurchase(r.Context(), e.OrderID, e.ID)
	if err != nil {
		log.Println("complete purchase:", err)
		http.Error(w, "retry", http.StatusInternalServerError) // Razorpay dobara bhejega
		return
	}

	// res == nil matlab pehle process ho chuka, email dobara nahi bhejna
	if res != nil {
		go func() {
			link := h.cfg.FrontendURL + "/access/" + res.AccessToken
			body := "Hi " + res.BuyerName + ",\n\n" +
				"Payment successful! " + res.CourseTitle + " ke notes aur WhatsApp group yahan milenge:\n\n" +
				link + "\n\n" +
				"Ye link kisi ke saath share mat karna.\n\n- Slantbit"
			if err := h.mailer.Send(res.BuyerEmail, "Welcome to "+res.CourseTitle, body); err != nil {
				log.Println("send email:", err)
			}
		}()
	}

	w.WriteHeader(http.StatusOK)
}