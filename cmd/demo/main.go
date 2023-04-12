package main

import (
	"fmt"
	"github.com/pablogolobaro/secret"
	"log"
	"os"
)

func main() {
	v := secret.File("fake-key", ".secrets")

	err := v.Set("demo-key", "some-data")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	plain, err := v.Get("demo-key")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Plain: %s", plain)
}
