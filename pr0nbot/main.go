package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
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
	cmd := `curl 'https://scrolller.com/api/media' -H 'origin: https://scrolller.com' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: en-US,en;q=0.9' -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.79 Safari/537.36' -H 'content-type: application/json' -H 'accept: */*' -H 'referer: https://scrolller.com/nsfw' -H 'authority: scrolller.com' -H 'cookie: _ga=GA1.2.661111315.1529359396; _gid=GA1.2.489592511.1529359396; _gat=1' --data-binary '[[` + myrand + `,null,0,10]]' --compressed`
	return cmd
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func sendpr0n(s *discordgo.Session, m *discordgo.MessageCreate) {
	out, err := exec.Command("sh", "-c", getCmd()).Output()
	if err != nil {
		fmt.Println(err)
	}
	grosseString := string(out)
	re := regexp.MustCompile(`([-a-zA-Z0-9_\/:.]+\.(jpg|mp4|webm))`)
	arrayDeString := re.FindAllString(grosseString, -1)
	if len(arrayDeString) == 0 {
		sendpr0n(s, m)
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
		break
	}

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
}
