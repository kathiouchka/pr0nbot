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
	"time"

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

func init2() {
	rand.Seed(time.Now().UnixNano())
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
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

// func getELOaoe4(s *discordgo.Session, m *discordgo.MessageCreate) {

// 	type Payload struct {
// 		Region       string `json:"region"`
// 		Versus       string `json:"versus"`
// 		MatchType    string `json:"matchType"`
// 		TeamSize     string `json:"teamSize"`
// 		SearchPlayer string `json:"searchPlayer"`
// 		Page         int    `json:"page"`
// 		Count        int    `json:"count"`
// 	}

// 	type Results struct {
// 		Count int `json:"count"`
// 		Items []struct {
// 			GameID       string      `json:"gameId"`
// 			UserID       string      `json:"userId"`
// 			RlUserID     int         `json:"rlUserId"`
// 			UserName     string      `json:"userName"`
// 			AvatarURL    interface{} `json:"avatarUrl"`
// 			PlayerNumber interface{} `json:"playerNumber"`
// 			Elo          int         `json:"elo"`
// 			EloRating    int         `json:"eloRating"`
// 			Rank         int         `json:"rank"`
// 			Region       string      `json:"region"`
// 			Wins         int         `json:"wins"`
// 			WinPercent   float64     `json:"winPercent"`
// 			Losses       int         `json:"losses"`
// 			WinStreak    int         `json:"winStreak"`
// 		} `json:"items"`
// 	}

// 	bodyMessage := strings.Fields(m.Content)
// 	var playerName string
// 	var matchType string
// 	if len(bodyMessage) != 3 {
// 		s.ChannelMessageSend(m.ChannelID, "Wrong entry. Example : `.elo Kathiou 1v1`")
// 		return
// 	}
// 	playerName = bodyMessage[1]
// 	matchType = bodyMessage[2]

// 	data := Payload{
// 		Region:       "0",
// 		Versus:       "players",
// 		MatchType:    "unranked",
// 		TeamSize:     matchType,
// 		SearchPlayer: playerName,
// 		Page:         1,
// 		Count:        100,
// 	}
// 	payloadBytes, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	body := bytes.NewReader(payloadBytes)

// 	req, err := http.NewRequest("POST", "https://api.ageofempires.com/api/ageiv/Leaderboard", body)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	req.Header.Set("Authority", "api.ageofempires.com")
// 	req.Header.Set("Sec-Ch-Ua", "\";Not A Brand\";v=\"99\", \"Chromium\";v=\"94\"")
// 	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
// 	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
// 	req.Header.Set("Sec-Ch-Ua-Platform", "\"Linux\"")
// 	req.Header.Set("Origin", "https://www.ageofempires.com")
// 	req.Header.Set("Sec-Fetch-Site", "same-site")
// 	req.Header.Set("Sec-Fetch-Mode", "cors")
// 	req.Header.Set("Sec-Fetch-Dest", "empty")
// 	req.Header.Set("Referer", "https://www.ageofempires.com/")
// 	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode == http.StatusOK {
// 		var result Results
// 		jsonBody, _ := ioutil.ReadAll(resp.Body)
// 		err = json.Unmarshal(jsonBody, &result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		s.ChannelMessageSend(m.ChannelID, "Elo: `"+strconv.Itoa(result.Items[0].Elo)+"` | Rank: `"+strconv.Itoa(result.Items[0].Rank)+"th`"+" | Winrate%: `"+strconv.Itoa(int(result.Items[0].WinPercent))+"`"+"| Winstreak: `"+strconv.Itoa(int(result.Items[0].WinStreak))+"`")
// 	} else {
// 		s.ChannelMessageSend(m.ChannelID, "No data found.")
// 	}

// }

// global historic of 10 last message from discord bot
var historic [10]string
var historicCount int = 0

func addToHistory(s *discordgo.Session, m *discordgo.MessageCreate) {

	msgInfo, _ := s.Channel(m.ChannelID)
	msgID := msgInfo.LastMessageID
	historic[historicCount] = msgID
	historicCount++
	if historicCount == 10 {
		historicCount = 0
	}
}
func remFromHistory(s *discordgo.Session, m *discordgo.MessageCreate) {
	if historicCount == 0 {
		historicCount = 10
	}
	s.ChannelMessageDelete(m.ChannelID, historic[historicCount-1])
	historicCount--
}

func remAllFromHistory(s *discordgo.Session, m *discordgo.MessageCreate) {
	for i := 0; i < len(historic); i++ {

		if historicCount == 0 {
			historicCount = 10
		}
		s.ChannelMessageDelete(m.ChannelID, historic[historicCount-1])
		historicCount--
	}
}

func sendpr0n(s *discordgo.Session, m *discordgo.MessageCreate) {
	re := regexp.MustCompile(`([-a-zA-Z0-9_\/:.]+(360).(jpg))`)
	if strings.Contains(m.Content, "vid") {
		// re = regexp.MustCompile(`([-a-zA-Z0-9_\/:.]+.(mp4))`)
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
