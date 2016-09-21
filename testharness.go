package testharness

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/thingful/thingfulx"
	"time"
)

var Number = "0"

// var Fetcher thingfulx.Fetcher

func Init(s string) {

	var Number = s
	fmt.Printf("input = %s\n", Number)
}

func Register(builder thingfulx.FetcherBuilder) (*Harness, error) {
	fetcher, err := builder()
	if err != nil {
		return nil, err
	}

	return &Harness{
		fetcher: fetcher,
	}, nil
}

type Harness struct {
	fetcher thingfulx.Fetcher
}

func (h *Harness) Run() {

	fmt.Println("Running fetcher:", h.fetcher.Provider().UID)
	fmt.Println("Provider:\n")
	spew.Dump(h.fetcher.Provider())

	fmt.Println("getting URLs:\n")

	timeout := time.Duration(30) * time.Second
	client := thingfulx.NewClient("thingful", timeout)
	delay := time.Duration(30) * time.Second

	URLs, err := h.fetcher.URLS(client, delay)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Got %v URLS\n", len(URLs))

	for _, u := range URLs {
		fmt.Println(u)
	}

}
