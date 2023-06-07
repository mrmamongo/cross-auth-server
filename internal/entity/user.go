package entity

type User struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	TelegramUsername string `json:"telegram_username"`
}
