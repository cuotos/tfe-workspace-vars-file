package main

import (
	"context"
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

func getVariablesForWorkspace(tfeClient *tfe.Client, workspaceID string) ([]*tfe.Variable, error) {
	variables := []*tfe.Variable{}

	opts := tfe.ListOptions{
		PageSize: 20,
	}

	for {
		variablesListOptions := tfe.VariableListOptions{
			ListOptions: opts,
		}

		returnedVariables, err := tfeClient.Variables.List(context.Background(), workspaceID, variablesListOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to get variables for workspace %v: %w", workspaceID, err)
		}

		variables = append(variables, returnedVariables.Items...)

		if returnedVariables.NextPage == 0 {
			break
		}

		opts.PageNumber = returnedVariables.NextPage
	}

	return variables, nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	token := mustGetEnv("TF_TOKEN")
	config := &tfe.Config{
		Token: token,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	vars, err := getVariablesForWorkspace(client, "ws-GT4xCxcD8AfBEwRT")
	if err != nil {
		log.Fatal(err)
	}

	outputString := strings.Builder{}

	for _, v := range vars {
		outputString.Write([]byte(fmt.Sprintf("%v=\"%v\"\n", v.Key, v.Value)))
	}

	fmt.Println(outputString.String())
}
