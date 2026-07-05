package dto

type User struct {
	Id          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	FullName    string `json:"full_name"`
	Photo       string `json:"photo"`
	Role        string `json:"role"`
	Verified    bool   `json:"verified"`
}

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
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	FullName    string `json:"full_name"`
	Photo       string `json:"photo"`
	Role        string `json:"role"`
}

type CreatedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type InternalErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
