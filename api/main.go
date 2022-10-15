package api

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
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
	Logger   *logrus.Entry
}

func (api *API) getEndpoint() string {
	if api.Port == "" {
		return fmt.Sprintf("%s:8728", api.Address)
	} else {
		return fmt.Sprintf("%s:%s", api.Address, api.Port)
	}
}

func (api *API) dial() (*routeros.Client, error) {
	endpoint := api.getEndpoint()
	if api.UseTLS {
		return routeros.DialTLS(endpoint, api.Username, api.Password, nil)
	}
	return routeros.Dial(endpoint, api.Username, api.Password)
}

func (api *API) Run(command string) ([]*proto.Sentence, error) {
	client, err := api.dial()
	if err != nil {
		api.Logger.Error(err)
		return []*proto.Sentence{}, err
	}
	defer client.Close()

	if api.Async {
		client.Async()
	}

	result, err := client.RunArgs(strings.Split(command, " "))
	if err != nil {
		api.Logger.Error(err)
	}

	return result.Re, err
}
