package main

import (
	"context"

	"github.com/hashicorp/go-tfe"
)

type TFClient struct {
	c *tfe.Client
}

func (t *TFClient) GetVariablesForWorkspace(workspaceID string) ([]*tfe.Variable, error) {
	vars := []*tfe.Variable{}

	opts := tfe.ListOptions{
		PageSize: 10,
	}

	for {

		listOptions := tfe.VariableListOptions{
			ListOptions: opts,
		}

		vs, err := t.c.Variables.List(context.Background(), workspaceID, listOptions)
		if err != nil {
			return nil, err
		}

		vars = append(vars, vs.Items...)

		if vs.NextPage == 0 {
			break
		}

		opts.PageNumber = vs.NextPage
	}

	return vars, nil
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
