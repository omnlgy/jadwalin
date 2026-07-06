package dto

type RegisterEmployeeRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	Email       string `json:"email" binding:"required,email"`
	Address     string `json:"address" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	Photo       string `json:"photo"`
}

type VerifyUserRequest struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type RegisterOTPRequest struct {
	Phone  string `json:"phone" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}
