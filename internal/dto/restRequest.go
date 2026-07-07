package dto

type RegisterEmployeeRequest struct {
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
	Phone  string `json:"phone" binding:"required" example:"+6281234567890"`
	UserID string `json:"user_id" binding:"required" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
}
