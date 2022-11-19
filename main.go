package main

import (
	"bytes"
	"github.com/lmika/micro-upload/micropub"
	"log"
	"os"
)

func main() {
	bearerAuth := os.Getenv("BEARER_AUTH")
	client := micropub.New("https://micro.blog/micropub", bearerAuth)

	log.Printf("Uploading image...")

	f, err := os.ReadFile("some-file.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	upload, err := client.Upload("https://example.micro.blog/", bytes.NewReader(f), "some-file.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	post := `Playing around with the Micropub API. Working on a small command line app to upload photos and`

	log.Printf("Making a post...")
	_, err = client.Post("https://example.micro.blog/", post)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done!")
}
