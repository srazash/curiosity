package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Photo struct {
	Id      float64 `json:"id"`
	Img_src string  `json:"img_src"`
}

var client *http.Client

func main() {

	sol := 1

	if len(os.Args) > 1 {
		if solArgs, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Println("Did dumbdumb type in a string instead of a number???")
			log.Fatal(err)
		} else {
			sol = solArgs
		}
	}

	url := fmt.Sprintf("https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos?sol=%d&api_key=DEMO_KEY", sol)

	client = &http.Client{Timeout: time.Second * 10}

	getPhoto(url)

}

func getJson(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(target)

	return nil
}

func getPhoto(url string) {
	var photo Photo

	if err := getJson(url, &photo); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%.0f: %s", photo.Id, photo.Img_src)
	}
}
