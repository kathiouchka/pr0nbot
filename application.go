package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"

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

var version = flag.String("version", "1.0.0", "the version number of the bot")

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
	flag.Parse()
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

func findUrls(content string) []string {
	var re *regexp.Regexp
	if strings.Contains(content, "vid") {
		re = regexp.MustCompile(`https://[-a-zA-Z0-9]+.scrolller.com/[-a-zA-Z0-9]+.mp4`)
	} else {
		re = regexp.MustCompile(`[-a-zA-Z0-9_/:.]+(1080).(jpg)`)
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
		return urls
	}
	return nil
}

func sendpr0n(s *discordgo.Session, m *discordgo.MessageCreate, counter int) {

	if counter == 5 {
		s.ChannelMessageSend(m.ChannelID, "`Error: I cannot find an image right now.`")
		counter = 0
	}

	urls := findUrls(m.Content)
	if len(urls) == 0 {
		sendpr0n(s, m, counter)
	}

	urlToSend := urls[rand.Intn(len(urls))]

	if m.Content == ".pr0n vid --debug" || m.Content == ".pr0n --debug" {
		message := fmt.Sprintf("```%.1950s```", strings.Join(urls, "\n"))
		_, discordError := s.ChannelMessageSend(m.ChannelID, message)
		if discordError != nil {
			s.ChannelMessageSend(m.ChannelID, discordError.Error())
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, urlToSend)
		addToHistory(s, m)
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
	case ".pr0n help":
		s.ChannelMessageSend(m.ChannelID, "``` .pr0n help | .kathiou | .pr0n user | .pr0n | .pr0n vid```")
	case ".pr0nbot":
		s.ChannelMessageSend(m.ChannelID, "Hi, I'm a naughty bot that can help you with some basic tasks. Type `.pr0n help` to see a list of available commands.")
	case ".pr0n delete":
		if m.Author.Username == "Kathiou" {
			info, _ := s.Channel(m.ChannelID)
			s.ChannelMessageDelete(m.ChannelID, info.LastMessageID)
			remFromHistory(s, m)
		}
	case ".pr0n deleteAll":
		if m.Author.Username == "Kathiou" {
			remAllFromHistory(s, m)
		}
	case ".kathiou":
		s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/633980782175584256/673619354360741912/kat.gif")
	case ".pr0n":
		{
			channel, err := s.Channel(m.ChannelID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			if !channel.NSFW {
				sendpr0n(s, m, 0)
			}
		}
	case ".pr0n vid":
		{
			channel, err := s.Channel(m.ChannelID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			if !channel.NSFW {
				sendpr0n(s, m, 0)
			}
		}
	case ".pr0n vid --debug":
		if m.Author.Username == "Kathiou" {
			channel, err := s.Channel(m.ChannelID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			if !channel.NSFW {
				sendpr0n(s, m, 0)
			}
		} else {
			user, err := s.User("400752755775373312")
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			s.ChannelMessageSend(m.ChannelID, "T'es pas <@!"+user.ID+"> mon pote.")
		}
	case ".pr0n --debug":
		if m.Author.Username == "Kathiou" {
			channel, err := s.Channel(m.ChannelID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			if !channel.NSFW {
				sendpr0n(s, m, 0)
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "T'es pas Kathiou mon pote.")
		}
	case ".pr0n --version":
		s.ChannelMessageSend(m.ChannelID, "`"+*version+"`")
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

// func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.Author.ID == s.State.User.ID {
// 		return
// 	}

// 	// Compile regular expressions
// 	arabicRegex, _ := regexp.Compile("(?i)arabe")
// 	amineRegex, _ := regexp.Compile("(?i)amine")
// 	bobRegex, _ := regexp.Compile("(?i)kathioubob")

// 	// Create a map to store command handlers
// 	commandHandlers := map[string]func(*discordgo.Session, *discordgo.MessageCreate){
// 		".pr0n help":        handleHelp,
// 		".pr0nbot":          handlePr0nbot,
// 		".pr0n delete":      handleDelete,
// 		".pr0n deleteAll":   handleDeleteAll,
// 		".kathiou":          handleKathiou,
// 		".pr0n":             handlePr0n,
// 		".pr0n vid":         handlePr0nVid,
// 		".pr0n vid --debug": handlePr0nVidDebug,
// 		".pr0n --debug":     handlePr0nDebug,
// 		".pr0n --version":   handlePr0nVersion,
// 		// Add more command handlers here
// 	}

// 	// Check if the message is a command and call the corresponding handler function
// 	if handler, ok := commandHandlers[m.Content]; ok {
// 		handler(s, m)
// 	}
// 	// Handle regular expressions
// 	if arabicRegex.MatchString(m.Content) {
// 		s.ChannelMessageSend(m.ChannelID, "(Amine)")
// 	}
// 	if amineRegex.MatchString(m.Content) {
// 		s.ChannelMessageSend(m.ChannelID, "(rebeu)")
// 	}
// 	if bobRegex.MatchString(m.Content) {
// 		s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/458438504129757186/1010225494869946470/kathioubob.png")
// 	}
// }

// // Define the handler functions

// func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	s.ChannelMessageSend(m.ChannelID, "``` .pr0n help | .kathiou | .pr0n user | .pr0n | .pr0n vid```")
// }

// func handlePr0nbot(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	s.ChannelMessageSend(m.ChannelID, "Hi, I'm a naughty bot that can help you with some basic tasks. Type `.pr0n help` to see a list of available commands.")
// }

// func handleDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.Author.Username == "Kathiou" {
// 		info, _ := s.Channel(m.ChannelID)
// 		s.ChannelMessageDelete(m.ChannelID, info.LastMessageID)
// 		remFromHistory(s, m)
// 	}
// }

// func handleDeleteAll(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.Author.Username == "Kathiou" {
// 		remAllFromHistory(s, m)
// 	}
// }

// func handleKathiou(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/633980782175584256/673619354360741912/kat.gif")
// }

// func handlePr0n(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	channel, err := s.Channel(m.ChannelID)
// 	if err != nil {
// 		s.ChannelMessageSend(m.ChannelID, err.Error())
// 	}
// 	if !channel.NSFW {
// 		sendpr0n(s, m, 0)
// 	}
// }

// func handlePr0nVid(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	channel, err := s.Channel(m.ChannelID)
// 	if err != nil {
// 		s.ChannelMessageSend(m.ChannelID, err.Error())
// 	}
// 	if !channel.NSFW {
// 		sendpr0n(s, m, 0)
// 	}
// }

// func handlePr0nVidDebug(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.Author.Username == "Kathiou" {
// 		channel, err := s.Channel(m.ChannelID)
// 		if err != nil {
// 			s.ChannelMessageSend(m.ChannelID, err.Error())
// 		}
// 		if !channel.NSFW {
// 			sendpr0n(s, m, 0)
// 		}
// 	} else {
// 		user, err := s.User("400752755775373312")
// 		if err != nil {
// 			s.ChannelMessageSend(m.ChannelID, err.Error())
// 		}
// 		s.ChannelMessageSend(m.ChannelID, "T'es pas <@!"+user.ID+"> mon pote.")
// 	}
// }
// func handlePr0nDebug(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.Author.Username == "Kathiou" {
// 		channel, err := s.Channel(m.ChannelID)
// 		if err != nil {
// 			s.ChannelMessageSend(m.ChannelID, err.Error())
// 		}
// 		if !channel.NSFW {
// 			sendpr0n(s, m, 0)
// 		}
// 	} else {
// 		s.ChannelMessageSend(m.ChannelID, "T'es pas Kathiou mon pote.")
// 	}
// }

// func handlePr0nVersion(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	s.ChannelMessageSend(m.ChannelID, "`"+*version+"`")
// }
