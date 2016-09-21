package testharness

import (
	"fmt"
	"github.com/thingful/thingfulx"
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
	spew.Dump(h.fetcher.Provider())
}
