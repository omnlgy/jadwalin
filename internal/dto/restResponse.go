package dto

import "time"

type User struct {
	Id          string `json:"id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	PhoneNumber string `json:"phone_number" example:"+628****7890"`
	Email       string `json:"email" example:"john.doe@example.com"`
	Address     string `json:"address" example:"Jl. Merdeka No. 1"`
	FullName    string `json:"full_name" example:"John Doe"`
	Photo       string `json:"photo" example:"http://example.com/photo.jpg"`
	Role        string `json:"role" example:"staff"`
	Verified    bool   `json:"verified" example:"true"`
}

type Treatment struct {
	Id          string  `json:"id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Name        string  `json:"name" example:"Haircut"`
	Description string  `json:"description" example:"Standard haircut"`
	Duration    int     `json:"duration" example:"30"`
	Price       float64 `json:"price" example:"50.00"`
}

type StaffSkill struct {
	Id              string  `json:"id"`
	UserID          string  `json:"user_id"`
	UserFullName    string  `json:"user_full_name"`
	UserPhoneNumber string  `json:"user_phone_number"`
	TreatmentID     string  `json:"treatment_id"`
	TreatmentName   string  `json:"treatment_name"`
	TreatmentPrice  float64 `json:"treatment_price"`
}

type FieldError struct {
	Field   string `json:"field" example:"phone_number"`
	Message string `json:"message" example:"invalid phone number"`
}

type BadRequestResponse struct {
	Code    int          `json:"code" example:"400"`
	Message string       `json:"message" example:"bad request"`
	Errors  []FieldError `json:"errors"`
}

type RegisterStaffResponse struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	FullName    string `json:"full_name"`
	Photo       string `json:"photo"`
	Role        string `json:"role"`
}

type CreatedResponse struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"User created successfully"`
	Data    any    `json:"data"`
}

type InternalErrorResponse struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal server error"`
}

type SuccessResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"OTP sent successfully"`
	Data    any    `json:"data,omitempty"`
}

type Meta struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"50"`
	TotalPages int   `json:"total_pages" example:"5"`
}

type PaginatedResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success"`
	Data    any    `json:"data"`
	Meta    Meta   `json:"meta"`
}

type UnauthorizedResponse struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"unauthorized"`
}

type ForbiddenResponse struct {
	Code    int    `json:"code" example:"403"`
	Message string `json:"message" example:"forbidden"`
}

type Booking struct {
	Id          string `json:"id"`
	ClientID    string `json:"client_id"`
	StaffID     string `json:"staff_id"`
	TreatmentID string `json:"treatment_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status"`
}

type AvailabelSlot struct {
	StartTime time.Time `json:"start_time" example:"09:00"`
	EndTime   time.Time `json:"end_time" example:"10:00"`
}

type GetAvailabelSlotResponse struct {
	AvailableSlots []AvailabelSlot `json:"available_slots"`
}
