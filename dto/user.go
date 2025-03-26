package dto

type UserloginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserUpdateRequest struct {
	Email       *string `json:"email" `
	FirstName   *string `json:"fname"`
	LastName    *string `json:"lname"`
	PhoneNumber *string `json:"phone"`
	Address     *string `json:"address"`
	Password    *string `json:"password"`
}

type UserDataRequest struct {
	Email *string `json:"email"`
}

type UserChangePasswordRequest struct {
	Email       *string `json:"email"`
	Password    *string `json:"password"`
	NewPassword *string `json:"newpassword"`
}
