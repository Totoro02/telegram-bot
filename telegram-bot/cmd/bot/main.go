package main

//6
// 53:19
import (
	"log"

	"github.com/Torgyn/telegram-bot/pkg/repository"
	"github.com/Torgyn/telegram-bot/pkg/repository/boltdb"
	"github.com/Torgyn/telegram-bot/pkg/server"
	"github.com/Torgyn/telegram-bot/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6032300960:AAHaa6YvPRcmAs51_pP4RnRjCCfNMYFYFbY")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	pocketClient, err := pocket.NewClient("107761-529c60e257b6954d3a9b6a9")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)
	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")
	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https://t.me/pocketLinks_bot")

	go func() {
		telegramBot.Start()
	}()
	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}

}

func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessToken))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestToken))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
