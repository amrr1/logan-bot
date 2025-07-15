package main

import (
    "fmt"
    "log"
    "os"
    "syscall"
    "os/signal"

    "github.com/bwmarrin/discordgo"
    "github.com/joho/godotenv"
    "go-logan-bot/src/bot" // Import your bot package
)

func main() {
    godotenv.Load()

    token := os.Getenv("BOT_KEY")

    sess, err := discordgo.New("Bot " + token)
    if err != nil {
        log.Fatal(err)
    }

    sess.AddHandler(bot.MessageHandler) // Use the handler from the bot package

    sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

    err = sess.Open()
    if err != nil {
        log.Fatal("error opening connection: ", err)
    }
    defer sess.Close()

    fmt.Println("Bot is now running. Press CTRL+C to exit.")

    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc
}