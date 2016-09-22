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
	timeout := time.Duration(60) * time.Second
	client := thingfulx.NewClient("thingful", timeout)
	delay := time.Duration(10) * time.Second

	URLs, err := h.fetcher.URLS(client, delay)
	if err != nil {
		// panic(err)
		fmt.Printf("Error from URLs: %s\n", err) // we should log this

	}

	showSize := int(math.Min(3, float64(len(URLs))))
	fmt.Printf("URLs has %d entry, showing first %d:\n", len(URLs), showSize)

	showURLs := URLs[0:showSize]

	for _, u := range showURLs {
		fmt.Println(u)
	}
	fmt.Println("\n\n")

	//FETCH
	fmt.Println("FETCH:\n")
	fmt.Printf("Fetch %d entry, showing result of first %d:\n", len(URLs), showSize)

	ctx := context.Background()
	clientFetch := thingfulx.NewClient("thingful", timeout)
	timeProvider := thingfulx.NewMockTimeProvider(time.Now())

	for i := 0; i < len(URLs); i++ {
		if i < showSize {
			fmt.Printf("fetching:  %s\n", URLs[i])
		}

		things, err := h.fetcher.Fetch(ctx, URLs[i], clientFetch, timeProvider)
		if err != nil {
			// panic(err)
			fmt.Printf("Error from Fetch: %s\n", err.Error()) // we should log this

		}

		if i < showSize {
			spew.Dump(things)
			fmt.Println("\n")
		}
	}

}

func (h *Harness) RunFetch(urls []string) {
	fmt.Printf("########### Running Fetcher: %s ########### \n", h.fetcher.Provider().UID)
	fmt.Println("FETCH:\n")
	timeout := time.Duration(60) * time.Second
	ctx := context.Background()
	clientFetch := thingfulx.NewClient("thingful", timeout)
	timeProvider := thingfulx.NewMockTimeProvider(time.Now())

	for _, u := range urls {
		fmt.Printf("fetching:  %s\n", u)
		things, err := h.fetcher.Fetch(ctx, u, clientFetch, timeProvider)
		if err != nil {
			// panic(err)
			fmt.Printf("Error from Fetch: %s\n", err.Error()) // we should log this

		}
		spew.Dump(things)
		fmt.Println("\n")
	}
}
