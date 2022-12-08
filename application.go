package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"

	"io/ioutil"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

// global historic of 10 last message from discord bot
var historic []string

func addToHistory(s *discordgo.Session, m *discordgo.MessageCreate) {
	msgInfo, _ := s.Channel(m.ChannelID)
	msgID := msgInfo.LastMessageID

	// Append the new message ID to the end of the slice
	historic = append(historic, msgID)

	// If the slice has more than 10 elements, remove the first element
	if len(historic) > 10 {
		historic = historic[1:]
	}
}

func remFromHistory(s *discordgo.Session, m *discordgo.MessageCreate) {
	// If the slice is empty, do nothing
	if len(historic) == 0 {
		return
	}

	// Delete the last message in the slice
	s.ChannelMessageDelete(m.ChannelID, historic[len(historic)-1])

	// Remove the last element from the slice
	historic = historic[:len(historic)-1]
}

func remAllFromHistory(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Delete all messages in the slice
	for _, msgID := range historic {
		s.ChannelMessageDelete(m.ChannelID, msgID)
	}

	// Clear the slice
	historic = nil
}

func sendpr0n(s *discordgo.Session, m *discordgo.MessageCreate) {

	re := regexp.MustCompile(`[-a-zA-Z0-9_/:.]+(360).(jpg)`)
	if strings.Contains(m.Content, "vid") {
		re = regexp.MustCompile(`https://[-a-zA-Z0-9]+.scrolller.com/[-a-zA-Z0-9]+.mp4`)
	}

	body := strings.NewReader("{\"query\":\" query DiscoverSubredditsQuery( $filter: MediaFilter $limit: Int $iterator: String $hostsDown: [HostDisk] ) { discoverSubreddits( isNsfw: true filter: $filter limit: $limit iterator: $iterator ) { iterator items { __typename url title secondaryTitle description createdAt isNsfw subscribers isComplete itemCount videoCount pictureCount albumCount isFollowing children( limit: 2 iterator: null filter: null disabledHosts: $hostsDown ) { iterator items { __typename url title subredditTitle subredditUrl redditPath isNsfw albumUrl isFavorite mediaSources { url width height isOptimized } } } } } } \",\"variables\":{\"limit\":30,\"filter\":null,\"hostsDown\":[\"NANO\",\"PICO\"]},\"authorization\":null}")
	req, err := http.NewRequest("POST", "https://api.scrolller.com/api/v2/graphql", body)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authority", "api.scrolller.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36")
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://scrolller.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://scrolller.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		bodyString := string(bodyBytes)
		urls := re.FindAllString(bodyString, -1)
		if len(urls) == 0 {
			sendpr0n(s, m)
		}

		urlToSend := urls[rand.Intn(len(urls))]

		if m.Content == ".pr0n vid --debug" {
			s.ChannelMessageSend(m.ChannelID, "inside")
			_, discordError := s.ChannelMessageSend(m.ChannelID, bodyString)
			if discordError != nil {
				s.ChannelMessageSend(m.ChannelID, discordError.Error())
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, urlToSend)
			addToHistory(s, m)
		}

	} else {
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa(resp.StatusCode))
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Compile regular expressions
	arabicRegex, _ := regexp.Compile("(?i)arabe")
	amineRegex, _ := regexp.Compile("(?i)amine")
	bobRegex, _ := regexp.Compile("(?i)kathioubob")

	// Handle commands
	switch m.Content {
	case ".help":
		s.ChannelMessageSend(m.ChannelID, "``` .help | .vid | .kathiou | .user | .pr0n | .pr0n vid```")
	case ".user":
		s.ChannelMessageSend(m.ChannelID, "Hi, I'm a naughty bot that can help you with some basic tasks. Type `.help` to see a list of available commands.")
	case ".delete":
		if m.Author.Username == "Kathiou" {
			info, _ := s.Channel(m.ChannelID)
			s.ChannelMessageDelete(m.ChannelID, info.LastMessageID)
			remFromHistory(s, m)
		}
	case ".deleteAll":
		if m.Author.Username == "Kathiou" {
			remAllFromHistory(s, m)
		}
	case ".kathiou":
		s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/633980782175584256/673619354360741912/kat.gif")
	case ".pr0n":
		sendpr0n(s, m)
	case ".pr0n vid":
		sendpr0n(s, m)
	case ".pr0n vid --debug":
		if m.Author.Username == "Kathiou" {
			sendpr0n(s, m)
		} else {
			s.ChannelMessageSend(m.ChannelID, "T'es pas Kathiou mon pote.")
		}
	}

	// Handle regular expressions
	if arabicRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "(Amine)")
	}
	if amineRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "(rebeu)")
	}
	if bobRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/458438504129757186/1010225494869946470/kathioubob.png")
	}
}
