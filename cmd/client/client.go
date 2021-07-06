package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SmAlexAl/ws_client/internal/userPool"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type TokenResponse struct {
	Token string `json:"token"`
}

func main() {
	up := userPool.UserPool{}

	up.InitListTokenRequest()

	for j := 1; j <= 1; j++ {
		newUser := userPool.User{
			Name:      "Profile" + strconv.Itoa(j),
			TokenByte: up.GetRandomTokenByte(),
		}
		//		var jsonStr = []byte(`{
		//	"osVersion": "8.1.1",
		//	"model": "iPhone 8",
		//	"platform": "Iphone",
		//	"pushToken": "string",
		//	"locale": "ru",
		//
		//	"applicationPackageName": "com.millcroft.inapp.sandbox",
		//	"applicationVersion": "1.0.0",
		//	"idfa": "9FF4ACCE-AEBF-4393-A354-E1B1FBF00B91",
		//	"installId": "9FF4ACCE-AEBF-4393-A354-E1B1FBF00123",
		//
		//	"udid": "FF60EE70-1F11-4880-BDE0-F908F2B18F88"
		//}`)
		//		bytes.NewBuffer(jsonStr)
		req, _ := http.NewRequest("POST", "http://localhost:8080/token", bytes.NewBuffer(newUser.TokenByte))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, _ := client.Do(req)

		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)

		respToken := &TokenResponse{}

		decoder.Decode(respToken)

		newUser.Token = respToken.Token

		go listenWs(&newUser)
	}

	fmt.Scanln()
}

type Object map[string]interface{}

type Message struct {
	Type    string `json:"type"`
	Token   string `json:"token"`
	Message string `json:"message"`
	Pin     bool   `json:"pin"`
}

func listenWs(user *userPool.User) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/chat/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
	updateToken(user)

	message := &Message{}
	message.Token = user.Token
	message.Pin = false
	message.Message = user.Name + " 1"
	message.Type = "authorization"

	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println("json error: ", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		fmt.Println("write error: ", err)
		return
	}
	count := 1
	ticker := time.NewTicker(time.Second * 10)
	tickerToken := time.NewTicker(time.Second * 25)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			count++
			message := &Message{}
			message.Pin = false
			message.Message = user.Name + " " + t.String()
			message.Type = "message"

			b, err := json.Marshal(message)
			if err != nil {
				fmt.Println("json error: ", err)
				return
			}
			err = c.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				fmt.Println("write error: ", err)
				return
			}
		case <-tickerToken.C:
			updateToken(user)
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func updateToken(user *userPool.User) {
	req, _ := http.NewRequest("POST", "http://localhost:8080/token", bytes.NewBuffer(user.TokenByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	respToken := &TokenResponse{}

	decoder.Decode(respToken)
	spew.Dump(respToken)

	user.Token = respToken.Token
}
