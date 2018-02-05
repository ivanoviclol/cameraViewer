package main

import (
	"github.com/gorilla/websocket"

	"net/http"

	"log"

	"fmt"

	"encoding/base64"

	"time"

	"io/ioutil"
)

func main() {

	http.HandleFunc("/connws/", ConnWs)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	http.HandleFunc("/templates/websocket.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/websocket.js")
	})

	err := http.ListenAndServe(":80", nil)

	if err != nil {

		log.Fatal("ListenAndServe: ", err)

	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

func ConnWs(w http.ResponseWriter, r *http.Request) {

	var ws, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error in socket!!!: ", err)
		return
	}

	var img64 []byte

	res := map[string]interface{}{}

	for {

		if err = ws.ReadJSON(&res); err != nil {

			if err.Error() == "EOF" {

				return

			}

			// ErrShortWrite means a write accepted fewer bytes than requested then failed to return an explicit error.

			if err.Error() == "unexpected EOF" {

				return

			}

			fmt.Println("Read : " + err.Error())

			return

		}

		res["a"] = "a"

		log.Println(res)

		for {

			files, _ := ioutil.ReadDir("./images")

			for _, f := range files {

				img64, _ = ioutil.ReadFile("./images/" + f.Name())

				str := base64.StdEncoding.EncodeToString(img64)

				res["img64"] = str

				if err = ws.WriteJSON(&res); err != nil {

					fmt.Println("watch dir - Write : " + err.Error())

					//return

				}

				time.Sleep(50 * time.Millisecond)

			}

			time.Sleep(50 * time.Millisecond)

		}

	}

}
