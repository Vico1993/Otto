package database

type BannedUser struct {
	TelegramId int64
	FirstName  string
	LastName   string
	UserName   string
	Lang       string
	IsBot      bool
}

func NewBannedUser(
	telegramId int64,
	firstName string,
	lastName string,
	userName string,
	lang string,
	isBot bool,
) *BannedUser {
	return &BannedUser{
		TelegramId: telegramId,
		FirstName:  firstName,
		LastName:   lastName,
		UserName:   userName,
		Lang:       lang,
		IsBot:      isBot,
	}
}
