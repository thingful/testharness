package testharness

import (
	"fmt"
	"math"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/thingful/thingfulx"
	"golang.org/x/net/context"
)

var (
	URLsErrorCount  int = 0
	URLsCount       int = 0
	FetchErrorCount int = 0
	ThingsCount     int = 0
	URLsError       []error
	FetchError      []string
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
	URLsCount = len(URLs)
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
			// err.Error() += " from " + URLs[i]
			if i < showSize {
				fmt.Printf("## ERROR from Fetch: %s\n", err.Error()) // we should log this
			}
			FetchError = append(FetchError, URLs[i]+"\n"+err.Error())
			FetchErrorCount += 1
		}
		ThingsCount += len(things)

		if i < showSize {
			fmt.Printf("there are %d things but printing first one\n", len(things))
			spew.Dump(things[0])
			fmt.Println("\n")
		}
	}

	fmt.Printf("\n########### SUMMARY ###########\n")
	fmt.Printf("Total URLs errors = %d\n", URLsErrorCount)
	for _, u := range URLsError {
		fmt.Println(u)
	}
	fmt.Printf("\nTotal URLs = %d\n", URLsCount)

	fmt.Printf("\nTotal fetch errors = %d\n", FetchErrorCount)
	for _, u := range FetchError {
		fmt.Println(u)
		fmt.Println()
	}
	fmt.Printf("\nTotal things fetched = %d\n", ThingsCount)
	if URLsErrorCount == 0 && FetchErrorCount == 0 {
		fmt.Printf("\nEverything seems to be OK\n\n")
	} else {
		fmt.Printf("\nThere seems to be problems\n\n")
	}

}

func (h *Harness) RunFetch(urls []string) {
	ThingsCount = 0
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
			FetchError = append(FetchError, err.Error())
			FetchErrorCount += 1
		}
		spew.Dump(things)
		fmt.Println("\n")
	}

	fmt.Printf("SUMMARY:")

	fmt.Printf("\nTotal fetch errors = %d\n", FetchErrorCount)
	for _, u := range FetchError {
		fmt.Println(u)
		fmt.Println()
	}
	fmt.Printf("\nTotal things fetched = %d\n", ThingsCount)
	if URLsErrorCount == 0 && FetchErrorCount == 0 {
		fmt.Printf("\nEverything seems to be OK\n\n")
	} else {
		fmt.Printf("\nThere seems to be problems\n\n")
	}

}
