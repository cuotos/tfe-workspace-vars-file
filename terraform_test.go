package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func startTestTFApiServer(t *testing.T) *url.URL {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	serverUrl, _ := url.Parse(server.URL)

	t.Cleanup(func() {
		server.Close()
	})

	return serverUrl
}

func getClientWithMockServer(t *testing.T) *TFClient {
	cfg := &tfe.Config{
		Token:   "NonNilString",
		Address: startTestTFApiServer(t).String(),
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create test client, this should never happen and is not a test failure. %v", err)
	}

	return client
}

func TestCanGetVarsFromWorkspace(t *testing.T) {

	inputVars := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	inputWorkspaceId := "testWorkspaceID"

	called := false

	client := getClientWithMockServer(t)
	client.c.Variables = MockVariables{
		ListFunc: func(workspaceID string, opts tfe.VariableListOptions) *tfe.VariableList {
			called = true
			assert.Equal(t, inputWorkspaceId, workspaceID)
			return &tfe.VariableList{
				Items:      createItemsList(inputVars),
				Pagination: &tfe.Pagination{NextPage: 0},
			}
		},
	}

	returnedVars, _ := client.GetVariablesForWorkspace(inputWorkspaceId)

	returnedVarsMap := map[string]string{}

	for _, v := range returnedVars {
		returnedVarsMap[v.Key] = v.Value
	}

	require.True(t, called)
	assert.Equal(t, fmt.Sprint(inputVars), fmt.Sprint(returnedVarsMap))
}

func TestProvidedListOptions(t *testing.T) {

	calledCount := 0 // how many times list was called
	var providedOps tfe.VariableListOptions

	client := getClientWithMockServer(t)
	client.c.Variables = MockVariables{
		ListFunc: func(workspaceID string, opts tfe.VariableListOptions) *tfe.VariableList {
			providedOps = opts
			listResponse := &tfe.VariableList{
				Items:      []*tfe.Variable{},
				Pagination: &tfe.Pagination{},
			}

			// mocking the behaviour of the tfe response pagination
			// if the current pagenumber is 4, return next pag = 0 to indicate no more pages
			if opts.PageNumber == 4 {
				listResponse.Pagination.NextPage = 0
			} else {
				listResponse.Pagination.NextPage = opts.PageNumber + 1
			}

			calledCount++

			return listResponse
		},
	}

	_, _ = client.GetVariablesForWorkspace("workspace")

	// make sure the list function was called twice
	assert.Equal(t, 5, calledCount, "list api was not called the expected number of times as expected by the pagination responses")
	assert.Equal(t, 10, providedOps.PageSize, "expected a different PageSize in list request pagination")
}

func createItemsList(vars map[string]string) []*tfe.Variable {
	items := []*tfe.Variable{}

	for k, v := range vars {
		items = append(items, &tfe.Variable{
			Key:   k,
			Value: v,
		})
	}
	return items
}

type MockVariables struct {
	ListFunc func(workspaceID string, opts tfe.VariableListOptions) *tfe.VariableList
}

// List all the variables associated with the given workspace.
func (mv MockVariables) List(ctx context.Context, workspaceID string, options tfe.VariableListOptions) (*tfe.VariableList, error) {
	return mv.ListFunc(workspaceID, options), nil
}

// Create is used to create a new variable.
func (mv MockVariables) Create(ctx context.Context, workspaceID string, options tfe.VariableCreateOptions) (*tfe.Variable, error) {
	panic("not implemented") // TODO: Implement
}

// Read a variable by its ID.
func (mv MockVariables) Read(ctx context.Context, workspaceID string, variableID string) (*tfe.Variable, error) {
	panic("not implemented") // TODO: Implement
}

// Update values of an existing variable.
func (mv MockVariables) Update(ctx context.Context, workspaceID string, variableID string, options tfe.VariableUpdateOptions) (*tfe.Variable, error) {
	panic("not implemented") // TODO: Implement
}

// Delete a variable by its ID.
func (mv MockVariables) Delete(ctx context.Context, workspaceID string, variableID string) error {
	panic("not implemented") // TODO: Implement
}
