package api

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/routeros.v2"
	"gopkg.in/routeros.v2/proto"
)

type API struct {
	Address  string
	Port     string
	Username string
	Password string
	UseTLS   bool
	Async    bool
}

func (api *API) getEndpoint() string {
	if api.Port == "" {
		return fmt.Sprintf("%s:8728", api.Address)
	} else {
		return fmt.Sprintf("%s:%s", api.Address, api.Port)
	}
}

func (api *API) dial() (*routeros.Client, error) {
	var client *routeros.Client
	err := errors.New("")
	endpoint := api.getEndpoint()
	if api.UseTLS {
		client, err = routeros.DialTLS(endpoint, api.Username, api.Password, nil)
	} else {
		client, err = routeros.Dial(endpoint, api.Username, api.Password)
	}
	return client, err
}

func (api *API) Run(command string) ([]*proto.Sentence, error) {
	client, err := api.dial()
	if err != nil {
		return []*proto.Sentence{}, err
	}
	defer client.Close()

	if api.Async {
		client.Async()
	}

	result, err := client.RunArgs(strings.Split(command, " "))
	if err != nil {
		recover()
	}
	return result.Re, err
}
