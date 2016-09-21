package testharness

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/thingful/thingfulx"
	"golang.org/x/net/context"
	"math"
	"time"
)

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

func (h *Harness) RunAll() {

	fmt.Printf("########### Running Fetcher: %s ########### \n", h.fetcher.Provider().UID)
	fmt.Println("Provider:\n")
	spew.Dump(h.fetcher.Provider())
	fmt.Println("\n\n")

	/// URLS
	fmt.Println("URLS:\n")
	timeout := time.Duration(30) * time.Second
	client := thingfulx.NewClient("thingful", timeout)
	delay := time.Duration(30) * time.Second

	URLs, err := h.fetcher.URLS(client, delay)
	if err != nil {
		panic(err)
	}

	showSize := int(math.Min(3, float64(len(URLs))))
	fmt.Printf("URLs has %d entry but only showing first %d:\n", len(URLs), showSize)

	showURLs := URLs[0:showSize]

	for _, u := range showURLs {
		fmt.Println(u)
	}
	fmt.Println("\n\n")

	//FETCH
	fmt.Println("FETCH:\n")
	ctx := context.Background()
	clientFetch := thingfulx.NewClient("thingful", timeout)
	timeProvider := thingfulx.NewMockTimeProvider(time.Now())

	for _, u := range showURLs {
		fmt.Printf("fetching:  %s\n", u)
		things, err := h.fetcher.Fetch(ctx, u, clientFetch, timeProvider)
		if err != nil {
			panic(err)
		}
		spew.Dump(things)
		fmt.Println("\n")
	}

}

func (h *Harness) RunFetch(urls []string) {
	fmt.Printf("########### Running Fetcher: %s ########### \n", h.fetcher.Provider().UID)
	fmt.Println("FETCH:\n")
	timeout := time.Duration(30) * time.Second
	ctx := context.Background()
	clientFetch := thingfulx.NewClient("thingful", timeout)
	timeProvider := thingfulx.NewMockTimeProvider(time.Now())

	for _, u := range urls {
		fmt.Printf("fetching:  %s\n", u)
		things, err := h.fetcher.Fetch(ctx, u, clientFetch, timeProvider)
		if err != nil {
			panic(err)
		}
		spew.Dump(things)
		fmt.Println("\n")
	}
}
