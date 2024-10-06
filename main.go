package main

import (
	"log"

	"github.com/AmbitiousJun/live-server/internal/web"
)

func main() {
	if err := web.Listen(5666); err != nil {
		log.Fatal(err)
	}
}
