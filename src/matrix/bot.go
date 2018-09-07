package matrix

import (
	"botinterface"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"plugins"
	"strings"
	"time"

	"events"

	"github.com/pkg/errors"
)

var (
	logPrefix = "MatrixBot: "
)

// The Bot struct holds parameters related to the bot
type Bot struct {
	receiveMessageChan chan events.ReceiveMessage
	sendMessageChan    chan events.SendMessage
	commandChan        chan events.Command
	server             string
	token              string
	pollingDone        chan bool

	pollingInterval time.Duration

	knownPlugins []plugins.Plugin
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
	DeviceID    string `json:"device_id"`
}

func (b Bot) apiCall(path string, method string, body string) (r []byte, e error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, b.server+"/_matrix"+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	// req.Header.Add("Authorization", "Bot "+b.token)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(response.Body)
}

// GetReceiveMessageChannel returns the channel which is used to notify
// about received messages from the bot. For DiscordBot these messages
// can be normal channel messages, whispers
func (b Bot) GetReceiveMessageChannel() chan events.ReceiveMessage {
	return b.receiveMessageChan
}

// GetSendMessageChannel returns the channel which is used to
// send messages using the bot. For DiscordBot these messages
// can be normal channel messages, whispers
func (b Bot) GetSendMessageChannel() chan events.SendMessage {
	return b.sendMessageChan
}

// GetCommandChannel gives a channel to control the bot from
// a plugin
func (b Bot) GetCommandChannel() chan events.Command {
	return b.commandChan
}

func (b *Bot) handlePolling() {
	log.Println(logPrefix + "TODO: IMPLEMENT SERVER POLLING")
	err := b.sendRoomMessage(string("!cJQhJDXTxLzZeuoHzw:matrix.abyle.org"), string("Hello Matrix World!"))
	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) startBot(doneChannel chan struct{}) {
	defer close(doneChannel)
	// do some message polling or whatever until stopped
	tickChan := time.Tick(b.pollingInterval)

	for {
		select {
		case <-tickChan:
			b.handlePolling()
		case <-b.pollingDone:
			return
		}
	}
}

func (b *Bot) login(username string, password string) (string, error) {
	// get login server
	response, err := b.apiCall("/client/r0/login", "POST", `{"type":"m.login.password", "user":"`+username+`", "password":"`+password+`"}`)
	if err != nil {
		return "", errors.Wrap(err, "apiCall failed")
	}

	log.Println(string(response))

	var channelResponseData loginResponse
	if err := json.Unmarshal(response, &channelResponseData); err != nil {
		return "", errors.Wrap(err, "json unmarshal failed")
	}

	if len(channelResponseData.AccessToken) > 0 {
		return channelResponseData.AccessToken, nil
	}

	return string(""), errors.New("could not login")
}

// CreateMatrixBot creates a new instance of a DiscordBot
func CreateMatrixBot(server string, username string, password string, token string) (*Bot, error) {
	log.Printf(logPrefix + "MatrixBot is CREATING itself")
	b := Bot{server: server}
	if len(token) == 0 {
		token, err := b.login(username, password)
		if err != nil {
			return nil, err
		}
		b.token = token
	} else {
		// just use the provided access token
		b.token = token
	}

	b.connectToMatrixServer()

	b.pollingDone = make(chan bool)
	b.pollingInterval = 1000 * time.Millisecond

	b.receiveMessageChan = make(chan events.ReceiveMessage)
	b.sendMessageChan = make(chan events.SendMessage)
	b.commandChan = make(chan events.Command)

	return &b, nil
}

func (b *Bot) connectToMatrixServer() error {
	response, err := b.apiCall("/client/r0/join/!cJQhJDXTxLzZeuoHzw:matrix.abyle.org?access_token="+b.token, "POST", `{}`)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}
	log.Println(string(response))
	return nil
}

func (b *Bot) startSendChannelReceiver() {
	for sendMsg := range b.sendMessageChan {
		switch sendMsg.Type {
		case events.MESSAGE:
			err := b.sendRoomMessage(sendMsg.Ident, sendMsg.Content)
			if err != nil {
				log.Println(err)
			}
		case events.WHISPER:
			log.Println(logPrefix + "events.WHISPER not implemented")
		default:
		}
	}
}

func (b *Bot) startCommandChannelReceiver() {
	for cmd := range b.commandChan {
		switch cmd.Command {
		case string("DemoCommand"):
			log.Println(logPrefix + "Received DemoCommand with server name" + cmd.Payload)
		default:
			log.Println(logPrefix + "Received unhandeled command" + cmd.Command)
		}
	}
}

// Start the Matrix Bot
func (b *Bot) Start(doneChannel chan struct{}) {
	log.Println(logPrefix + "MatrixBot is STARTING")
	go b.startBot(doneChannel)
	go b.startSendChannelReceiver()
	go b.startCommandChannelReceiver()
}

// Stop the Matrix Bot
func (b Bot) Stop() {
	log.Println(logPrefix + "MatrixBot is SHUTING DOWN")

	b.pollingDone <- true

	defer close(b.receiveMessageChan)
}

// Status returns the current status of MatrixBot
func (b *Bot) Status() botinterface.BotStatus {
	return botinterface.BotStatus{Running: true}
}

// AddPlugin adds the give plugin to the current bot
func (b *Bot) AddPlugin(plugin plugins.Plugin) {
	plugin.ConnectChannels(b.GetReceiveMessageChannel(), b.GetSendMessageChannel(), b.GetCommandChannel())
	b.knownPlugins = append(b.knownPlugins, plugin)
}
