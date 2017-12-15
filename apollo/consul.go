package Apollo

import (
	consulapi "github.com/hashicorp/consul/api"
	"errors"
	"math/rand"
)

func GetServer(name string)  (*Server,error) {
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, err
	}
	services, err := client.Agent().Services()
	var newServices []*consulapi.AgentService
	for _, v := range services {
		if v.Service == name {
			newServices = append(newServices, v)
		}
	}
	count := len(newServices)
	if count == 0 {
		return nil, errors.New("an't found sercie " + name)
	} else {
		index := rand.Intn(count)
		v := newServices[index]
		return &Server{v.Port, v.Address}, nil
	}
}
