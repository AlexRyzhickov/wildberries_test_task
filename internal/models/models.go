package models

type UserGrade struct {
	UserId        string `json:"user_id" validate:"required"`
	PostpaidLimit int    `json:"postpaid_limit,omitempty"`
	Spp           int    `json:"spp,omitempty"`
	ShippingFee   int    `json:"shipping_fee,omitempty"`
	ReturnFee     int    `json:"return_fee,omitempty"`
}
