package subscriber

type Message struct {
	PhoneNumber string `json:"phone_number"`
	Id          int    `json:"id"`
	Amount      int    `json:"amount"`
}
