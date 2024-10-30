package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/carlmjohnson/requests"
	"github.com/gofiber/fiber/v2"

	// "github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/template/html/v2"
)

var store *Storage 

func RenderForm(c *fiber.Ctx) error {
	
    return c.Render("login", store)
}

func ProcessSettings(c *fiber.Ctx) error {
	if store == nil {
		return c.Render("login", fiber.Map{})
	}

	settingsData := getSettings()

	store.Settings = *settingsData

    return c.Render("form", fiber.Map{"Store": *store, "Id":  store.Id, "Token":  store.Token})
}

func ProcessStateInstance(c *fiber.Ctx) error {
	if store == nil {
		return c.Render("login", fiber.Map{})
	}

	stateInst := getStateInstance()

	store.StateInst = stateInst

    return c.Render("form", fiber.Map{"Store": *store, "Id":  store.Id, "Token":  store.Token})
}

func ProcessMessage(c *fiber.Ctx) error {
	if store == nil {
		return c.Render("login", fiber.Map{})
	}

	message := c.FormValue("message")

	textMessage := TextMessage{ChatID: store.Settings.Wid, Message: message}
	
	messageId := sendMessage(textMessage,  store.Id,  store.Token)

	store.LastMessage = messageId
    return c.Render("form", fiber.Map{"Store": *store, "Id":  store.Id, "Token":  store.Token})
}
func ProcessFileUrl(c *fiber.Ctx) error {
	if store == nil {
		return c.Render("login", fiber.Map{})
	}

	caption := c.FormValue("caption")
	url := c.FormValue("url")

    fileName := regexp.MustCompile(`[^/]+\.[a-z]+$`).FindString(url)

	store.LastMessage = sendFileByUrl(FileMessage{URLFile: url, FileName: fileName, ChatID: store.Settings.Wid, Caption: caption})

	

    return c.Render("form", fiber.Map{"Store": *store, "Id":  store.Id, "Token":  store.Token})
}
func Login(c *fiber.Ctx) error {
	if store == nil {
		store = &Storage{}
	}

	id := c.FormValue("id")
	token := c.FormValue("token")
	store.Id =fmt.Sprint(id)
	store.Token = fmt.Sprint(token)
	
	settingsData := getSettings()
	if (settingsData == nil) {
		return  c.Render("login", fiber.Map{})
	}
	store.Settings = *settingsData
	store.LastMessage = ResultSendMessage{}
	
	return c.Render("form", fiber.Map{"Store": *store, "Id": id, "Token": token})
}
func getSettings() *ResultSettings {
	var result ResultSettings
	
	err := requests.
		URL(fmt.Sprintf("https://1103.api.green-api.com/waInstance%s/getSettings/%s", store.Id, store.Token)).
		ToJSON(&result).
		Fetch(context.Background())

	if err != nil {
		fmt.Println("could not connect to example.com:", err)
		return nil
	}

	return &result
}
func getStateInstance() ResultStateInstance {
	var result ResultStateInstance

	err := requests.
		URL(fmt.Sprintf("https://1103.api.green-api.com/waInstance%s/getStateInstance/%s", store.Id, store.Token)).
		ToJSON(&result).
		Fetch(context.Background())

	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	
	return result
}
func sendMessage(message TextMessage, id string, token string) ResultSendMessage {
	var result ResultSendMessage
	
	err := requests.
		URL(fmt.Sprintf("https://1103.api.green-api.com/waInstance%s/sendMessage/%s", id, token)).
		BodyJSON(&message).
		ToJSON(&result).
		Fetch(context.Background())

	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	
	return result
}
func sendFileByUrl(fileMessage FileMessage) ResultSendMessage {
	var result ResultSendMessage

	err := requests.
		URL(fmt.Sprintf("https://1103.api.green-api.com/waInstance%s/sendFileByUrl/%s", store.Id, store.Token)).
		BodyJSON(&fileMessage).
		ToJSON(&result).
		Fetch(context.Background())

	if err != nil {
		fmt.Println("could not connect to example.com:", err)

	}
	
	return result
}

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	app.Static("/", "./static")
	

	app.Get("/", RenderForm)
	app.Post("/login", Login)
	app.Post("/message", ProcessMessage)
	app.Post("/fileUrl", ProcessFileUrl)
	app.Get("/settings", ProcessSettings)
	app.Get("/stateInstance", ProcessStateInstance)


	app.Listen("10.100.102.6:3000")
}

type ResultSettings struct {
    
    CountryInstance              string `json:"countryInstance"`
    TypeAccount                  string `json:"typeAccount"`
    WebhookUrl                   string `json:"webhookUrl"`
    WebhookUrlToken              string `json:"webhookUrlToken"`
    DelaySendMessagesMilliseconds int    `json:"delaySendMessagesMilliseconds"`
    MarkIncomingMessagesReaded   string `json:"markIncomingMessagesReaded"`
    MarkIncomingMessagesReadedOnReply string `json:"markIncomingMessagesReadedOnReply"`
    SharedSession                string `json:"sharedSession"`
    OutgoingWebhook              string `json:"outgoingWebhook"`
    OutgoingMessageWebhook       string `json:"outgoingMessageWebhook"`
    OutgoingAPIMessageWebhook    string `json:"outgoingAPIMessageWebhook"`
    IncomingWebhook              string `json:"incomingWebhook"`
    DeviceWebhook                string `json:"deviceWebhook"`
    StatusInstanceWebhook        string `json:"statusInstanceWebhook"`
    StateWebhook                 string `json:"stateWebhook"`
    EnableMessagesHistory        string `json:"enableMessagesHistory"`
    KeepOnlineStatus             string `json:"keepOnlineStatus"`
    PollMessageWebhook           string `json:"pollMessageWebhook"`
    IncomingBlockWebhook         string `json:"incomingBlockWebhook"`
    IncomingCallWebhook          string `json:"incomingCallWebhook"`
	Wid                          string `json:"wid"`
}

type ResultStateInstance struct {
    StateInstance   string `json:"stateInstance"`
   
}
type ResultSendMessage struct {
    IdMessage       string `json:"idMessage"`
   
}

type FileMessage struct {
    ChatID          string  `json:"chatId"`
    URLFile         string  `json:"urlFile"`
    FileName        string  `json:"fileName"`
    Caption         string  `json:"caption"`
    QuotedMessageID *string `json:"quotedMessageId,omitempty"`
}
type TextMessage struct {
    ChatID 			string  `json:"chatId"`
    Message         string  `json:"message"`
    QuotedMessageID *string `json:"quotedMessageId,omitempty"`
    LinkPreview     *bool   `json:"linkPreview,omitempty"`
}

type Storage struct {
	
    Settings		ResultSettings `json:"settings"`
	StateInst		ResultStateInstance  `json:"stateInst"`
    LastMessage			ResultSendMessage `json:"message"`
	Id 				string`json:"id"`
	Token 			string`json:"token"`
	
	
}