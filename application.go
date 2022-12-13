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

type scrolllerStruct struct {
	Data struct {
		GetSubreddit struct {
			Children struct {
				Iterator string `json:"iterator"`
				Items    []struct {
					Typename         string      `json:"__typename"`
					ID               int         `json:"id"`
					URL              string      `json:"url"`
					Title            string      `json:"title"`
					SubredditID      int         `json:"subredditId"`
					SubredditTitle   string      `json:"subredditTitle"`
					SubredditURL     string      `json:"subredditUrl"`
					RedditPath       string      `json:"redditPath"`
					IsNsfw           bool        `json:"isNsfw"`
					AlbumURL         interface{} `json:"albumUrl"`
					HasAudio         interface{} `json:"hasAudio"`
					FullLengthSource interface{} `json:"fullLengthSource"`
					GfycatSource     interface{} `json:"gfycatSource"`
					RedgifsSource    interface{} `json:"redgifsSource"`
					OwnerAvatar      interface{} `json:"ownerAvatar"`
					Username         interface{} `json:"username"`
					DisplayName      interface{} `json:"displayName"`
					IsPaid           interface{} `json:"isPaid"`
					Tags             interface{} `json:"tags"`
					IsFavorite       bool        `json:"isFavorite"`
					MediaSources     []struct {
						URL         string `json:"url"`
						Width       int    `json:"width"`
						Height      int    `json:"height"`
						IsOptimized bool   `json:"isOptimized"`
					} `json:"mediaSources"`
					BlurredMediaSources interface{} `json:"blurredMediaSources"`
				} `json:"items"`
			} `json:"children"`
		} `json:"getSubreddit"`
	} `json:"data"`
}

var version = flag.String("version", "v.0.0.TEST", "the version number of the bot")

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
	var body *strings.Reader
	if strings.Contains(content, "vid") {
		re = regexp.MustCompile(`https://[-a-zA-Z0-9]+.scrolller.com/[-a-zA-Z0-9]+.mp4`)
	} else {
		re = regexp.MustCompile(`[-a-zA-Z0-9_/:.]+(1080).(jpg)`)
	}
	// body := strings.NewReader("{\"query\":\" query DiscoverSubredditsQuery( $filter: MediaFilter $limit: Int $iterator: String $hostsDown: [HostDisk] ) { discoverSubreddits( isNsfw: true filter: $filter limit: $limit iterator: $iterator ) { iterator items { __typename url title secondaryTitle description createdAt isNsfw subscribers isComplete itemCount videoCount pictureCount albumCount isFollowing children( limit: 2 iterator: null filter: null disabledHosts: $hostsDown ) { iterator items { __typename url title subredditTitle subredditUrl redditPath isNsfw albumUrl isFavorite mediaSources { url width height isOptimized } } } } } } \",\"variables\":{\"limit\":30,\"filter\":null,\"hostsDown\":[\"NANO\",\"PICO\"]},\"authorization\":null}")
	index := strings.Index(content, " ")
	subreddit := strings.ToLower(content[index+1:])
	if subreddit == ".pr0n" || subreddit == "vid" {
		// .pr0n
		fmt.Println("SUBREDDIT = ", subreddit)
		body = strings.NewReader("{\"query\":\" query DiscoverSubredditsQuery( $filter: MediaFilter $limit: Int $iterator: String $hostsDown: [HostDisk] ) { discoverSubreddits( isNsfw: true filter: $filter limit: $limit iterator: $iterator ) { iterator items { __typename url title secondaryTitle description createdAt isNsfw subscribers isComplete itemCount videoCount pictureCount albumCount isFollowing children( limit: 2 iterator: null filter: null disabledHosts: $hostsDown ) { iterator items { __typename url title subredditTitle subredditUrl redditPath isNsfw albumUrl isFavorite mediaSources { url width height isOptimized } } } } } } \",\"variables\":{\"limit\":30,\"filter\":null,\"hostsDown\":[\"NANO\",\"PICO\"]},\"authorization\":null}")
	} else {
		fmt.Println("SUBREDDIT ELSE = ", subreddit)
		body = strings.NewReader(`{"query":" query SubredditQuery( $url: String! $filter: SubredditPostFilter $iterator: String ) { getSubreddit(url: $url) { children( limit: 50 iterator: $iterator filter: $filter disabledHosts: null ) { iterator items { __typename id url title subredditId subredditTitle subredditUrl redditPath isNsfw albumUrl hasAudio fullLengthSource gfycatSource redgifsSource ownerAvatar username displayName isPaid tags isFavorite mediaSources { url width height isOptimized } blurredMediaSources { url width height isOptimized } } } } } ","variables":{"url":"/r/` + subreddit + `","filter":null,"hostsDown":null},"authorization":null}`)
	}
	req, err := http.NewRequest("POST", "https://api.scrolller.com/api/v2/graphql", body)
	fmt.Println(req)
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
		fmt.Println(bodyString)
		urls := re.FindAllString(bodyString, -1)
		return urls
	}
	return nil
}

func sendpr0n(s *discordgo.Session, m *discordgo.MessageCreate, counter int) {
	urls := findUrls(m.Content)
	if len(urls) == 0 && counter < 3 {
		counter++
		sendpr0n(s, m, counter)
	}

	if counter == 0 {
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
	} else if counter == 2 {
		s.ChannelMessageSend(m.ChannelID, "`Error: no img`")
	}
}

// func getData(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
// 	if err != nil {
// 		s.ChannelMessageSend(m.ChannelID, err.Error())
// 	}
// 	var botMessages int
// 	for _, message := range messages {
// 		if message.Author.ID == s.State.User.ID {
// 			botMessages++
// 		}
// 	}
// 	if botMessages > 0 {
// 		s.ChannelMessageSend(m.ChannelID, strconv.Itoa(botMessages))
// 	}
// }

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Compile regular expressions
	arabicRegex, _ := regexp.Compile("(?i)arabe")
	amineRegex, _ := regexp.Compile("(?i)amine")
	bobRegex, _ := regexp.Compile("(?i)kathioubob")
	subredditRegexp := regexp.MustCompile(`^\.pr0n\s+\w+`)

	// Handle commands
	switch m.Content {
	case ".pr0n help":
		s.ChannelMessageSend(m.ChannelID, "``` .pr0n help | .kathiou | .pr0n | .pr0n vid```")
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
			if channel.NSFW {
				sendpr0n(s, m, 0)
			} else {
				s.ChannelMessageSend(m.ChannelID, "This channel is not NSFW!")
			}
		}
	case ".pr0n rand":
		number := rand.Intn(100)
		if number == 0 {
			number += 1
		}
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa(number))
	case ".pr0n vid":
		{
			fmt.Println("BBBB")
			channel, err := s.Channel(m.ChannelID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			if channel.NSFW {
				sendpr0n(s, m, 0)
			} else {
				s.ChannelMessageSend(m.ChannelID, "This channel is not NSFW!")
			}
		}
	case subredditRegexp.FindString(m.Content):
		{
			fmt.Println("AAAAA")
			sendpr0n(s, m, 0)
		}
	// case ".pr0n vid --debug":
	// 	if m.Author.Username == "Kathiou" {
	// 		channel, err := s.Channel(m.ChannelID)
	// 		if err != nil {
	// 			s.ChannelMessageSend(m.ChannelID, err.Error())
	// 		}
	// 		if channel.NSFW {
	// 			sendpr0n(s, m, 0)
	// 		} else {
	// 			s.ChannelMessageSend(m.ChannelID, "This channel is not NSFW!")
	// 		}
	// 	} else {
	// 		user, err := s.User("400752755775373312")
	// 		if err != nil {
	// 			s.ChannelMessageSend(m.ChannelID, err.Error())
	// 		}
	// 		s.ChannelMessageSend(m.ChannelID, "T'es pas <@!"+user.ID+"> mon pote.")
	// 	}
	// case ".pr0n --debug":
	// 	if m.Author.Username == "Kathiou" {
	// 		channel, err := s.Channel(m.ChannelID)
	// 		if err != nil {
	// 			s.ChannelMessageSend(m.ChannelID, err.Error())
	// 		}
	// 		if channel.NSFW {
	// 			sendpr0n(s, m, 0)
	// 		} else {
	// 			s.ChannelMessageSend(m.ChannelID, "This channel is not NSFW!")
	// 		}
	// 	} else {
	// 		s.ChannelMessageSend(m.ChannelID, "T'es pas Kathiou mon pote.")
	// 	}
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
