package testharness

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/thingful/thingfulx"
	"golang.org/x/net/context"
	"math"
	"time"
)

var (
	URLsErrorCount  int = 0
	FetchErrorCount int = 0
	URLsError       []error
	FetchError      []error
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
		fmt.Printf("## ERROR from URLs: %s\n", err.Error()) // we should log this
		URLsError = append(URLsError, err)
		URLsErrorCount += 1
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
			fmt.Printf("Fetching:  %s\n", URLs[i])
		}

		things, err := h.fetcher.Fetch(ctx, URLs[i], clientFetch, timeProvider)
		if err != nil {
			// panic(err)
			fmt.Printf("## ERROR from Fetch: %s\n", err.Error()) // we should log this
			FetchError = append(FetchError, err)
			FetchErrorCount += 1
		}

		if i < showSize {
			spew.Dump(things)
			fmt.Println("\n")
		}
	}

	fmt.Printf("\n########### SUMMARY ###########\n")
	fmt.Printf("URLsErrorCount = %d\n", URLsErrorCount)
	for _, u := range URLsError {
		fmt.Println(u)
	}
	fmt.Printf("\nFetchErrorCount = %d\n", FetchErrorCount)
	for _, u := range FetchError {
		fmt.Println(u)
	}
	if URLsErrorCount == 0 && FetchErrorCount == 0 {
		fmt.Printf("\nEverything seems to be OK\n\n")
	} else {
		fmt.Printf("\nThere seems to be problems\n\n")
	}
}

func (h *Harness) RunFetch(urls []string) {
	FetchErrorCount = 0
	FetchError = FetchError[:0]
	fmt.Printf("########### Running Fetcher: %s ########### \n", h.fetcher.Provider().UID)
	fmt.Println("FETCH:\n")
	timeout := time.Duration(60) * time.Second
	ctx := context.Background()
	clientFetch := thingfulx.NewClient("thingful", timeout)
	timeProvider := thingfulx.NewMockTimeProvider(time.Now())

	for _, u := range urls {
		fmt.Printf("Fetching:  %s\n", u)
		things, err := h.fetcher.Fetch(ctx, u, clientFetch, timeProvider)
		if err != nil {
			// panic(err)
			fmt.Printf("## ERROR from Fetch: %s\n", err.Error()) // we should log this
			FetchError = append(FetchError, err)
			FetchErrorCount += 1
		}
		spew.Dump(things)
		fmt.Println("\n")
	}

	fmt.Printf("SUMMARY:")
	fmt.Printf("FetchErrorCount = %d\n", FetchErrorCount)
	for _, u := range FetchError {
		fmt.Println(u)
	}
}
