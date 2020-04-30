package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"

	// "strings"
	"bytes"
	"io/ioutil"
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

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func getCmd() string {

	myrand := strconv.Itoa(random(1, 2500))
	cmd := `curl 'https://scrolller.com/api/media' -H 'origin: https://scrolller.com' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: en-US,en;q=0.9' -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.79 Safari/537.36' -H 'content-type: application/json' -H 'accept: */*' -H 'referer: https://scrolller.com/nsfw' -H 'authority: scrolller.com' -H 'cookie: _ga=GA1.2.661111315.1529359396; _gid=GA1.2.489592511.1529359396; _gat=1' --data-binary '[[` + myrand + `,null,0,2]]' --compressed`
	return cmd
}

func sendpr0n(s *discordgo.Session, m *discordgo.MessageCreate) {

	re := regexp.MustCompile(`([-a-zA-Z0-9_\/:.]+\.(jpg|mp4|webm))`)
	chanName := regexp.MustCompile(`"(\w+)"`)
	myrand := strconv.Itoa(random(1, 6000))
	var jsonStr = []byte(`[[` + myrand + `,null,0,2]]`)
	req, err := http.NewRequest("POST", "https://scrolller.com/api/media", bytes.NewBuffer(jsonStr))
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "scrolller.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.122 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://scrolller.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://scrolller.com/nsfw")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7")
	req.Header.Set("Cookie", "__cfduid=d282e4ea05bcbe9d31e8693b0901d45a51587755050; _ga=GA1.2.317970967.1587755062; _gid=GA1.2.1812126065.1588158281; _gat=1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		bodyString := string(bodyBytes)
		images := re.FindAllString(bodyString, -1)
		if len(images) == 0 {
			// restart func until we find an image
			sendpr0n(s, m)
		}
		s.ChannelMessageSend(m.ChannelID, "Subreddit `"+chanName.FindString(bodyString)+"`")
		s.ChannelMessageSend(m.ChannelID, "https://scrolller.com/media/"+images[0])
	}

	// grosseString := string(out)
	// re := regexp.MustCompile(`([-a-zA-Z0-9_\/:.]+\.(jpg|mp4|webm))`)
	// arrayDeString := re.FindAllString(grosseString, -1)
	// if len(arrayDeString) == 0 {
	// 	sendpr0n(s, m)
	// }
	// vals := arrayDeString
	// r := rand.New(rand.NewSource(time.Now().Unix()))
	// for _, i := range r.Perm(len(vals)) {
	// 	val := vals[i]
	// 	gfy := val
	// 	reg := regexp.MustCompile(`gfy`)
	// 	anus := reg.FindString(gfy)
	// 	if anus == "gfy" {
	// 		s.ChannelMessageSend(m.ChannelID, val)
	// 		break
	// 	}
	// 	s.ChannelMessageSend(m.ChannelID, "https://scrolller.com/media/"+val)
	// 	break
	// }

}

func sendpr0nbomb(s *discordgo.Session, m *discordgo.MessageCreate) {

	out, err := exec.Command("sh", "-c", getCmd()).Output()
	if err != nil {
		fmt.Println(err)
	}
	grosseString := string(out)
	re := regexp.MustCompile(`([-a-zA-Z0-9_\/:.]+\.(jpg|mp4|webm))`)
	arrayDeString := re.FindAllString(grosseString, -1)
	if len(arrayDeString) == 0 {
		sendpr0nbomb(s, m)
	}

	vals := arrayDeString
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(vals)) {
		val := vals[i]
		gfy := val
		reg := regexp.MustCompile(`gfy`)
		anus := reg.FindString(gfy)
		if anus == "gfy" {
			s.ChannelMessageSend(m.ChannelID, val)
			break
		}
		s.ChannelMessageSend(m.ChannelID, "https://scrolller.com/media/"+val)
		if i == 5 {
			break
		}
	}

}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == ".pr0n" {
		sendpr0n(s, m)
	}
	if m.Content == ".pr0nbomb" {
		sendpr0nbomb(s, m)
	}
	if m.Content == ".ocl0cks" {
		s.ChannelMessageSend(m.ChannelID, "https://slacker.scrolller.com/media/38da5d.mp4")
	}

}
