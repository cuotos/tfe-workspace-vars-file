package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-tfe"
)

var (
	workspaceID string
)

func mustGetEnv(v string) string {
	value := os.Getenv(v)
	if value == "" {
		log.Fatalf("Required env variable %v was not found", v)
	}
	return value
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	flag.StringVar(&workspaceID, "w", "", "workspace to generate vars file from")
	flag.Parse()

	if workspaceID == "" {
		return errors.New("no workspace provided")
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	token := mustGetEnv("TF_TOKEN")
	config := &tfe.Config{
		Token: token,
	}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	vars, err := client.GetVariablesForWorkspace(workspaceID)
	if err != nil {
		return err
	}

	outputString := strings.Builder{}

	for _, v := range vars {
		outputString.Write([]byte(fmt.Sprintf("%v=\"%v\"\n", v.Key, v.Value)))
	}

	fmt.Println(outputString.String())

	return nil
}
