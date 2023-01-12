package main

import (
    "gopkg.in/telegram-bot-api.v4"
    "log"
)

const botToken = "Put your bot token here!"

func main() {
    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true
    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi")
        //msg.ReplyToMessageID = update.Message.MessageID
        btn := tgbotapi.KeyboardButton{
            RequestLocation: true,
            Text: "Gimme where u live!!",
        }
        msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})
        bot.Send(msg)
    }
}



