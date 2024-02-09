package dtos

type UseCoupon struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Code        string `json:"code"`
}
