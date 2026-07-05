package dto

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type BadRequestResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

type RegisterEmployeeResponse struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	FullName    string `json:"full_name"`
	Photo       string `json:"photo"`
	Role        string `json:"role"`
}
