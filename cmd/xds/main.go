package main

import (
	"log"

	"github.com/cognicraft/xds"
)

func main() {
	s := xds.New(":1883")
	log.Fatal(s.Run())
}
