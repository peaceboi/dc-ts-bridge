package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
  Token string
)

const KuteGoAPIURL = "https://kutego-api-xxxxx-ew.a.run.app"

func init () {
  flag.StringVar(&Token, "t", "", "Bot Token")
  flag.Parse()
}

func main() {

  //new discord session using the provided bot token
  dg, err  := discordgo.new("Bot" + Token)
  if err != nil {
      fmt.Println("Error creating discord session,", err)
      return
  }

  //register messageCreate function as a callback for MessageCreate events
  dg.AddHandler(messageCreate)

  //example - we only care about receiving message events
  dg.Identify.Intents = discordgo.IntentsGuildMessages

  //open websocket connection to discord and begin listening
  err = dg.Open()
  if err != nil {
    fmt.Println("error opening connection,", err)
    return
  }

  //wait here until CTRL+C or other termination signal is received
  fmt.Println("Bot is now running. Press CTRL+C to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill,)
  <-sc

  //cleanly close down discord session
  dg.Close()

  type hänger struct {
    Name string `json: "name"`
  }

  /*this function will be called every time a message is created on any channel the bot has access to (due to AddHandler above)*/
  func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate)

    //ignore all messages from the bot itself
    if m.Author.ID == s.State.User.ID {
      return
    }

    if m.Content == "/hänger" {
      
    }
}