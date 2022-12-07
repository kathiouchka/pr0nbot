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
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
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
	re := regexp.MustCompile(`([-a-zA-Z0-9_\/:.]+(360).(jpg))`)
	if strings.Contains(m.Content, "vid") {
		re = regexp.MustCompile(`(https://[-a-zA-Z0-9]+.scrolller.com/[-a-zA-Z0-9]+.mp4)`)
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
		fmt.Println(urls)
		if len(urls) == 0 {
			sendpr0n(s, m)
		}

		urlToSend := urls[rand.Intn(len(urls))]
		s.ChannelMessageSend(m.ChannelID, urlToSend)
		addToHistory(s, m)

	} else {
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa(resp.StatusCode))
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == ".pr0n" || m.Content == ".pr0n vid" {
		sendpr0n(s, m)
	}
	if m.Content == ".helpzer" {
		s.ChannelMessageSend(m.ChannelID, "``` .pr0n | .pr0n vid | .kathi0u | .helpzer | kathioubob ```")
	}
	regexDeRebeu, _ := regexp.Compile("(?i)arabe")
	if regexDeRebeu.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "(Amine)")
	}
	regexDAmine, _ := regexp.Compile("(?i)amine")
	if regexDAmine.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "(rebeu)")
	}
	regexBOB, _ := regexp.Compile("(?i)kathioubob")
	if regexBOB.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/458438504129757186/1010225494869946470/kathioubob.png")
	}

	if m.Author.Username == "Kathiou" && m.Content == ".delete" {
		info, _ := s.Channel(m.ChannelID)
		s.ChannelMessageDelete(m.ChannelID, info.LastMessageID)
		remFromHistory(s, m)
	}
	if m.Author.Username == "Kathiou" && m.Content == ".deleteAll" {
		remAllFromHistory(s, m)
	}

	if m.Content == ".kathi0u" {
		s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/633980782175584256/673619354360741912/kat.gif")
	}
}