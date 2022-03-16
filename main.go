package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-tfe"
)

func mustGetEnv(v string) string {
	value := os.Getenv(v)
	if value == "" {
		log.Fatalf("Required env variable %v was not found", v)
	}
	return value
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	token := mustGetEnv("TF_TOKEN")
	config := &tfe.Config{
		Token: token,
	}

	client, err := NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	vars, err := client.GetVariablesForWorkspace("ws-GT4xCxcD8AfBEwRT")
	if err != nil {
		log.Fatal(err)
	}

	outputString := strings.Builder{}

	for _, v := range vars {
		outputString.Write([]byte(fmt.Sprintf("%v=\"%v\"\n", v.Key, v.Value)))
	}

	fmt.Println(outputString.String())
}
