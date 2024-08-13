package telegram

import (
	// "fmt"
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	commandStart = "start"

	replyStartTemplayte = "Hi! For saving links you should autarization: "

	replyAlreadyAuthorization = "Hi! You are authorized! You can send links!"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}

}
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Link is saved!")
	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		msg.Text = "It's invalid link"
		_, err = b.bot.Send(msg)
		return err
	}
	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		msg.Text = "You are not authorized! Use command /start"
		_, err = b.bot.Send(msg)
		return err
	}

	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		msg.Text = "I cann't save link :( Try again Letter!"
		_, err = b.bot.Send(msg)
		return err
	}
	_, err = b.bot.Send(msg)
	return err

}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, replyAlreadyAuthorization)
	b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "I don't know this command :(")
	_, err := b.bot.Send(msg)
	return err
}
