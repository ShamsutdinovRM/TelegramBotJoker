package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	botToken := "your token"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = respond(botUrl, update)
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return restResponse.Result, nil
}

func getJoke(str string) (string, error) {
	resp, err := http.Get(str)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var chuck Response
	err = json.Unmarshal(body, &chuck)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	return chuck.Value, nil
}

func random() string {
	r, _ := getJoke("https://api.chucknorris.io/jokes/random")

	return r
}

func respond(botUrl string, update Update) error {
	var BotMessage BotMessage
	if update.Message.Text == "/random" {
		BotMessage.ChatId = update.Message.Chat.ChatId
		BotMessage.Text = random()
		buf, err := json.Marshal(BotMessage)
		if err != nil {
			return err
		}
		_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
		return nil
	} else {
		BotMessage.ChatId = update.Message.Chat.ChatId
		BotMessage.Text = "I don't understand your message, please send /random, to get a joke!"
		buf, err := json.Marshal(BotMessage)
		if err != nil {
			return err
		}
		_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
		return nil
	}
}
