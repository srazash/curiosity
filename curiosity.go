package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Images struct {
	Photos []struct {
		ImgSrc string `json:"img_src"`
	} `json:"photos"`
}

var client *http.Client

func main() {

	sol := 1

	if len(os.Args) > 1 {
		if solArgs, err := strconv.Atoi(os.Args[1]); err != nil {
			log.Fatal(err)
		} else {
			sol = solArgs
		}
	}

	url := fmt.Sprintf("https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos?sol=%d&api_key=DEMO_KEY", sol)

	client = &http.Client{Timeout: time.Second * 10}

	fmt.Printf("ｃｕｒｉｏｓｉｔｙ░ｖｅｒｓ．░０．９９　（くごフ）\n【ｓｏｆｔｗａｒｅ　ｂｙ　ｒｓｈ】\n...starting download of ｃｕｒｉｏｓｉｔｙ 【SOL %d】 img files...\n", sol)

	getPhotos(sol, url)

}

func getPhotos(sol int, url string) {
	var photos Images
	dir := fmt.Sprintf("sol %s", strconv.Itoa(sol))

	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		log.Fatal(err)
	}

	if err := getJson(url, &photos); err != nil {
		log.Fatal(err)
	} else {
		for i, p := range photos.Photos {
			resp, err := http.Get(p.ImgSrc)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			filename := strings.Split(p.ImgSrc, "/")
			file := fmt.Sprintf("%s/%s", dir, filename[len(filename)-1])
			perca := float64(i) + 1
			percb := float64(len(photos.Photos))
			perc := perca / percb * 100
			fmt.Printf("【%.2f%%】\tdownloading file %d of %d...\n", perc, i+1, len(photos.Photos))

			out, err := os.Create(file)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Fatal(err)
			}

		}

	}

	fmt.Printf("look for 【%s】 dir in the PWD...", dir)

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
