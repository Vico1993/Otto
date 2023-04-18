package database

type User struct {
	ChatId     int64
	TelegramId int64
	FirstName  string
	LastName   string
	UserName   string
	Lang       string
	IsBot      bool
	IsBanned   bool
}

func NewUser(
	chatId int64,
	telegramId int64,
	firstName string,
	lastName string,
	userName string,
	lang string,
	isBot bool,
	isBanned bool,
) *User {
	return &User{
		ChatId:     chatId,
		TelegramId: telegramId,
		FirstName:  firstName,
		LastName:   lastName,
		UserName:   userName,
		Lang:       lang,
		IsBot:      isBot,
		IsBanned:   isBanned,
	}
}
