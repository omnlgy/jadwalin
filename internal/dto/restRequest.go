package dto

type RegisterUserRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164" example:"+6281234567890"`
	Email       string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Address     string `json:"address" binding:"required" example:"Jl. Merdeka No. 1"`
	FullName    string `json:"full_name" binding:"required" example:"John Doe"`
	Photo       string `json:"photo" example:"http://example.com/photo.jpg"`
}

type VerifyUserRequest struct {
	Phone string `json:"phone" binding:"required" example:"+6281234567890"`
	OTP   string `json:"otp" binding:"required" example:"123456"`
}

type RegisterOTPRequest struct {
	Phone  string `json:"phone" binding:"required" example:"+628****7890"`
	UserID string `json:"user_id" binding:"required" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
}

type LoginRequest struct {
	Phone string `json:"phone" binding:"required,e164" example:"+628****7890"`
}

type UpdateUserRequest struct {
	PhoneNumber string `json:"phone_number" example:"+628****7890"`
	Email       string `json:"email" example:"john.doe@example.com"`
	Address     string `json:"address" example:"Jl. Merdeka No. 1"`
	FullName    string `json:"full_name" example:"John Doe"`
	Photo       string `json:"photo" example:"http://example.com/photo.jpg"`
}

type CreateTreatmentRequest struct {
	Name        string  `json:"name" binding:"required" example:"Haircut"`
	Description string  `json:"description" example:"Standard haircut"`
	Duration    int     `json:"duration" binding:"required,min=1" example:"30"`
	Price       float64 `json:"price" binding:"required,min=0" example:"50.00"`
}

type UpdateTreatmentRequest struct {
	Name        string  `json:"name" example:"Haircut"`
	Description string  `json:"description" example:"Standard haircut"`
	Duration    int     `json:"duration" binding:"omitempty,min=1" example:"30"`
	Price       float64 `json:"price" binding:"omitempty,min=0" example:"50.00"`
}

type AssignSkillRequest struct {
	UserID      string `json:"user_id" binding:"required" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	TreatmentID string `json:"treatment_id" binding:"required" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
}
