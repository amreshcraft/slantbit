package models

// Public listing: whatsapp link aur pdf key kabhi JSON mein nahi jayenge
type Course struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	PricePaise  int64  `json:"pricePaise"`

	WhatsappLink string `json:"-"`
	PdfKey       string `json:"-"`
}

type Access struct {
	EnrollmentID int64
	CourseID     int64
	CourseTitle  string
	WhatsappLink string
	PdfKey       string
}

type PurchaseResult struct {
	BuyerName   string
	BuyerEmail  string
	CourseTitle string
	AccessToken string
}