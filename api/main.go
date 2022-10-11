package api

import (
	"log"
	"strings"

	"gopkg.in/routeros.v2"
	"gopkg.in/routeros.v2/proto"
)

type API struct {
	Address  string
	Username string
	Password string
	UseTLS   bool
	Async    bool
}

func (api *API) dial() (*routeros.Client, error) {
	if api.UseTLS {
		return routeros.DialTLS(api.Address, api.Username, api.Password, nil)
	}
	return routeros.Dial(api.Address, api.Username, api.Password)
}

func (api *API) Run(command string) ([]*proto.Sentence, error) {
	client, err := api.dial()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if api.Async {
		client.Async()
	}

	result, err := client.RunArgs(strings.Split(command, " "))
	if err != nil {
		log.Fatal(err)
	}

	return result.Re, err
}
