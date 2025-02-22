package internal

import (
	"fmt"
	"strings"

	"github.com/go-routeros/routeros/v3"
	"github.com/go-routeros/routeros/v3/proto"
)

type Api struct {
	Address  string
	Port     string
	Username string
	Password string
	UseTLS   bool
	Async    bool
}

func (api *Api) getEndpoint() string {
	if api.Port == "" {
		return fmt.Sprintf("%s:8728", api.Address)
	} else {
		return fmt.Sprintf("%s:%s", api.Address, api.Port)
	}
}

func (api *Api) dial() (*routeros.Client, error) {
	endpoint := api.getEndpoint()
	if api.UseTLS {
		return routeros.DialTLS(endpoint, api.Username, api.Password, nil)
	} else {
		return routeros.Dial(endpoint, api.Username, api.Password)
	}
}

func (api *Api) Run(command string) ([]*proto.Sentence, error) {
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
		return []*proto.Sentence{}, err
	}
	return result.Re, err
}
