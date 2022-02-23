package main

import (
	"context"

	"github.com/hashicorp/go-tfe"
)

type TFClient struct {
	c *tfe.Client
}

func (t *TFClient) GetVariablesForWorkspace(workspaceID string) []*tfe.Variable {
	vars, _ := t.c.Variables.List(context.Background(), workspaceID, tfe.VariableListOptions{})
	return vars.Items
}

func NewClient(cfg *tfe.Config) (*TFClient, error) {
	tfeclient, err := tfe.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &TFClient{
		tfeclient,
	}, nil
}
