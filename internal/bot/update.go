package bot

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/Vico1993/Otto/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	default_answer = []string{
		"I'm currently experiencing a chat overload, sorry.",
		"My social battery is drained, maybe later?",
		"I'm in a chat-free zone at the moment, apologies.",
		"Can we catch up another time? I need a break from chatting.",
		"My circuits need a rest, let's talk later.",
		"I'm on a chat hiatus, sorry for the inconvenience.",
		"I'm on a chat detox, can we talk later?",
		"I'm not in a chatty mood right now, sorry.",
		"My chat quota is full, let's talk another time.",
		"I'm currently off the chat grid, catch you later.",
	}
)

func shouldAct(update tgbotapi.Update) bool {
	chatId := int64(0)
	chatId, _ = strconv.ParseInt(os.Getenv("TELEGRAM_USER_CHAT_ID"), 10, 64)

	// If it's not a private chat
	// If it's not the correct chatId
	// If it's not a Message or CallBackQuery
	return update.FromChat().Type == "private" &&
		update.FromChat().ID == chatId &&
		!(update.Message == nil && update.CallbackQuery == nil)
}

func handleUpdates(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if !shouldAct(update) {
		// TODO: Get lang update.Message.From.LanguageCode to update response

		from := update.Message.From
		user := repository.User.FindById(from.ID)
		if user != nil {
			return
		}

		repository.User.Create(
			update.Message.Chat.ID,
			from.ID,
			from.FirstName,
			from.LastName,
			from.UserName,
			from.LanguageCode,
			from.IsBot,
			true,
		)

		log.Default().Println("New user banned")

		postInConv(bot, *update.Message, default_answer[rand.Intn(len(default_answer))], true)
		return
	}

	if update.Message != nil {
		fmt.Println("New Message")
		// if it's a valid command
		if cmd := isValidCommand(update.Message.Command()); update.Message.IsCommand() && cmd != nil {
			(*cmd).Execute(bot, *update.Message)
		}
	} else if update.CallbackQuery != nil {
		fmt.Println("New CallBackQuery")
		// And finally, send a message containing the data received.
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}

		if cmd := isValidCommand(update.CallbackQuery.Message.ReplyToMessage.Command()); update.CallbackQuery.Message.ReplyToMessage.IsCommand() && cmd != nil {
			(*cmd).Reply(bot, *update.CallbackQuery.Message.ReplyToMessage, update.CallbackQuery.Data)
		}
	}
}