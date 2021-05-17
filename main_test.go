package main

import (
	"fmt"
	"testing"

	"github.com/vahidmostofi/minaria/sdk/client"
	"github.com/vahidmostofi/minaria/sdk/client/heath"
)

func TestHealthCheck(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)
	r, err := c.Heath.CheckHealthStatus(heath.NewCheckHealthStatusParams())

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(r)
}
