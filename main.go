package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	witai "github.com/wit-ai/wit-go/v2"
)

var (
	Token string
)

const KuteGoAPIURL = "http://192.168.1.57:8080"

/* func init() {
	flag.StringVar(&Token, "t", "", "Bot token")
	flag.Parse()
}
*/

func main() {

	godotenv.Load(".env")
	Token := os.Getenv("DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session, ", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection ", err)
		return
	}

	fmt.Println("Bot is now running. CTLR-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

type Gopher struct {
	Name string `json:"name"`
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// ignore my own messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!commands" {

		commandlist :=
			`
		Bot commands 
		!gophers - send a list of current gophers*
		!gopher - sends the dr who gopher*
		!random - sends a random gopher*

		!hello - Sends test message - Hello back to user

		!ask - sends a question to wolfram alpha

		* Only if API is running
		
		`

		_, err := s.ChannelMessageSend(m.ChannelID, "\n\n"+commandlist+"\n")
		if err != nil {
			fmt.Println(err)
		}
	}

	if m.Content == "!hello" {

		_, err := s.ChannelMessageSend(m.ChannelID, "Hi there.")
		if err != nil {
			fmt.Println(err)
		}

	}

	if m.Content == "!gopher" {

		response, err := http.Get(KuteGoAPIURL + "/gopher/" + "dr-who")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, "dr-who.png", response.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Cannot get dr-who gopher! :(")
		}
	}

	if m.Content == "!random" {

		response, err := http.Get(KuteGoAPIURL + "/gopher/random")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, "random-gopher.png", response.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Cannot get random gopher! :(")
		}
	}

	if m.Content == "!gophers" {

		response, err := http.Get(KuteGoAPIURL + "/gophers/")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}

			var data []Gopher
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println(err)
			}

			var gophers strings.Builder
			for _, gopher := range data {
				gophers.WriteString(gopher.Name + "\n")
			}

			_, err = s.ChannelMessageSend(m.ChannelID, gophers.String())
			if err != nil {
				fmt.Println(err)
			}

		} else {
			fmt.Println("Error: Cannot get list of gophers! :(")
		}
	}

	if m.Content == "!ask" {
		fmt.Println("Hello")

		client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))

		client.Parse(&witai.MessageRequest{s.})

	}

	if m.Content == "!lookup" {
		fmt.Println("Lookup")
	}

}
