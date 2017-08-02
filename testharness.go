package testharness

import (
	"fmt"
	"math"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/temoto/robotstxt"
	"github.com/thingful/thingfulx"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	// urlsErrorCount indicates how many errors occured during URLs method
	urlsErrorCount int

	// urlsCount indicates number of URLs returned by URLs methog
	urlsCount int

	// fetchErrorCount indicates how many errors occured during Fetch method
	fetchErrorCount int

	// thingsCount indicates how many things successfully fetched from Fetch method
	thingsCount int

	// emptyThingsCount indicates how many URLs return without error but also without thing
	emptyThingsCount int

	// urlsError stores errors from URLs method
	urlsError []error

	// fetchError stores errors from Fetch method
	fetchError []string

	// whiteListed indicated if this indexer is whitelisted from checking robots.txt
	whiteListed bool
)

// Register function register the indexer to test harness
// Init by passing the indexer and whitelisted bool to Register method
func Register(builder thingfulx.IndexerBuilder, whitelisted bool) (*Harness, error) {
	whiteListed = whitelisted
	indexer, err := builder()
	if err != nil {
		return nil, err
	}

	return &Harness{
		indexer: indexer,
	}, nil
}

// Harness contain instance of Indexer, and ability to run Indexer
type Harness struct {
	indexer thingfulx.Indexer
}

// RunAll runs URL method of this indexer first
// Then use the result URLs to fetch
// Then for each URL to fetch, get one dataURL from first thing returned to access
func (h *Harness) RunAll(ctx context.Context, fetchInterval time.Duration, totalFetch int) {

	fmt.Printf("########### Running Indexer: %s ########### \n", h.indexer.UID())

	/// ############## URLS ###############
	fmt.Println("URLS:\n")
	timeout := time.Duration(60) * time.Second
	client := thingfulx.NewClient("thingful", timeout)
	delay := time.Duration(10) * time.Second

	URLs, err := h.indexer.URLS(ctx, client, delay)
	if err != nil {
		fmt.Printf("## ERROR from URLs: %s\n", err.Error()) // we should log this
		urlsError = append(urlsError, err)
		urlsErrorCount++
	}

	showSize := min([]int{3, len(URLs)})
	urlsCount = len(URLs)
	fmt.Printf("URLs has %d entry, showing first %d:\n", len(URLs), showSize)

	showURLs := URLs[0:showSize]

	for _, u := range showURLs {
		fmt.Println(u)
	}
	fmt.Println("\n\n")
	// WE CHECK IF THESE URLS ARE ALLOWED HERE ###############
	fmt.Printf("CHECKING FOR ROBOTS.TXT FOR ALL URLS\n")
	allAllowed, allowErr := checkURLs(URLs)
	if allowErr != nil {
		fmt.Printf("there is error from checkURLs %s\n", allowErr)
		return
	}

	if !allAllowed {
		fmt.Printf("the URLs are blocked by robots.txt\n")
		return
	}

	// ############### FETCH ###############
	totalFetch = min([]int{totalFetch, len(URLs)})
	showSize = min([]int{totalFetch, 3})
	fmt.Println("FETCH:\n")
	fmt.Printf("Fetch %d entry, showing result of first %d:\n", totalFetch, showSize)

	clientFetch := thingfulx.NewClient("thingful", timeout)
	timeProvider := thingfulx.NewMockTimeProvider(time.Now())
	var dataURLs []string // dataURL to check

	for i := 0; i < totalFetch; i++ {

		if i < showSize {
			fmt.Printf("Fetching:  %s\n", URLs[i])
		}

		bytes, err := h.indexer.Fetch(ctx, URLs[i], clientFetch)

		if err != nil {
			if i < showSize {
				fmt.Printf("## ERROR from Fetch: %s\n", err.Error()) // we should log this
			}
			fetchError = append(fetchError, URLs[i]+"\n"+err.Error())
			fetchErrorCount++
		}

		things, err := h.indexer.Parse(bytes, URLs[i], timeProvider)

		if err != nil {
			if i < showSize {
				fmt.Printf("## ERROR from Fetch: %s\n", err.Error()) // we should log this
			}
			fetchError = append(fetchError, URLs[i]+"\n"+err.Error())
			fetchErrorCount++
		} else {

			if len(things) == 0 {
				emptyThingsCount++
			} else {
				thingsCount += len(things)
				dataURLs = append(dataURLs, things[0].Endpoint.URL) // save 1 dataUrl from each fetch to test
			}

		}

		if i < showSize {
			if len(things) > 0 {
				fmt.Printf("there are %d things but printing first one\n", len(things))
				spew.Dump(things[0])
				fmt.Println("\n")
			} else {
				fmt.Printf("things slice is empty\n")
			}

		}

		time.Sleep(fetchInterval)
	}

	// ############### ACCESS ###############
	fmt.Printf("########### Checking Access for: %s ########### \n", h.indexer.UID())
	successAccessCount := 0
	failureAccessCount := 0
	foundUniqueURL := 0 // to store how many have we found
	for _, u := range dataURLs {

		foundUniqueURL = 0

		fmt.Printf("Accessing:  %s\n", u)
		bytes, err := h.indexer.Fetch(ctx, u, clientFetch)
		if err != nil {
			fmt.Printf("## ERROR from Fetch: %s\n", err.Error())
			fetchError = append(fetchError, err.Error())
			fetchErrorCount++
		}

		things, err := h.indexer.Parse(bytes, u, timeProvider)
		if err != nil {
			fmt.Printf("## ERROR from Parse: %s\n", err.Error())
			fetchError = append(fetchError, err.Error())
			fetchErrorCount++
		} else {

			if len(things) == 0 {
				fmt.Printf("ERROR this URL: %s returns nothing\n", u)
				failureAccessCount++
			} else {
				for i := 0; i < len(things); i++ {
					if u == things[i].Endpoint.URL { // check if one of "things" contain the same urls that used to access it
						foundUniqueURL++
						fmt.Printf("Found same unique dataURL: %s \n", things[i].Endpoint.URL)
					}
				}
			}

		}

		if foundUniqueURL == 1 { // we only expect to find 1 copy of dataURL
			fmt.Printf("SUCCESS found one match: \n")
			successAccessCount++
		} else if foundUniqueURL > 1 { // if we found more than one, give warning
			fmt.Printf("ERROR found %d of the same unique URL: \n", foundUniqueURL)
			failureAccessCount++
		} else if foundUniqueURL == 0 { // if we found nothing, also give warning
			fmt.Printf("ERROR can't find anything that match: %s \n", u)
			failureAccessCount++
		}

	}

	fmt.Printf("\n########### SUMMARY ###########\n")
	fmt.Printf("Total URLs errors = %d\n", urlsErrorCount)
	for _, u := range urlsError {
		fmt.Println(u)
	}
	fmt.Printf("\nTotal URLs = %d\n", urlsCount)

	fmt.Printf("\nTotal fetch errors = %d\n", fetchErrorCount)
	for _, u := range fetchError {
		fmt.Println(u)
		fmt.Println()
	}
	fmt.Printf("\nTotal things fetched = %d\n", thingsCount)
	fmt.Printf("\nTotal empty things = %d\n", emptyThingsCount)

	fmt.Printf("\nTotal things access attemp = %d\n", len(dataURLs))
	fmt.Printf("\nTotal things access successfully = %d\n", successAccessCount)

	if urlsErrorCount == 0 && urlsCount > 0 && fetchErrorCount == 0 && failureAccessCount == 0 {
		fmt.Printf("\nEverything seems to be OK\n\n")
	} else {
		fmt.Printf("\nThere seems to be problems\n\n")
	}

}

// RunFetch `fetch` the specified URLs then `parse` the content
func (h *Harness) RunFetch(ctx context.Context, urls []string, fetchInterval time.Duration) {
	thingsCount = 0
	emptyThingsCount = 0
	fetchErrorCount = 0
	fetchError = fetchError[:0]
	fmt.Printf("########### Running Indexer: %s ########### \n", h.indexer.UID())
	fmt.Println("FETCH:\n")
	fmt.Printf("CHECKING FOR ROBOTS.TXT FOR ALL URLS\n")
	allAllowed, allowErr := checkURLs(urls)
	if allowErr != nil {
		fmt.Printf("there is error from checkURLs %s\n", allowErr)
		return
	}

	if !allAllowed {
		fmt.Printf("the URLs are blocked by robots.txt\n")
		return
	}

	timeout := time.Duration(60) * time.Second
	clientFetch := thingfulx.NewClient("thingful", timeout)
	timeProvider := thingfulx.NewMockTimeProvider(time.Now())

	for _, u := range urls {
		fmt.Printf("Fetching:  %s\n", u)
		bytes, err := h.indexer.Fetch(ctx, u, clientFetch)

		if err != nil {
			fmt.Printf("## ERROR from Fetch: %s\n", err.Error())
			fetchError = append(fetchError, err.Error())
			fetchErrorCount++
		}

		things, err := h.indexer.Parse(bytes, u, timeProvider)

		if err != nil {
			fmt.Printf("## ERROR from Parse: %s\n", err.Error())
			fetchError = append(fetchError, err.Error())
			fetchErrorCount++
		} else {

			if len(things) == 0 {
				emptyThingsCount++
			} else {
				thingsCount += len(things)
			}
		}
		spew.Dump(things)
		fmt.Println("\n")
		time.Sleep(fetchInterval)
	}

	fmt.Printf("SUMMARY:")

	fmt.Printf("\nTotal fetch errors = %d\n", fetchErrorCount)
	for _, u := range fetchError {
		fmt.Println(u)
		fmt.Println()
	}
	fmt.Printf("\nTotal things fetched = %d\n", thingsCount)
	fmt.Printf("\nTotal empty things = %d\n", emptyThingsCount)
	if urlsErrorCount == 0 && fetchErrorCount == 0 {
		fmt.Printf("\nEverything seems to be OK\n\n")
	} else {
		fmt.Printf("\nThere seems to be problems\n\n")
	}

}

//RunAccess fetch(and parse) the dataURLs, then check that ONLY one thing returned from Fetch has the same dataURL
func (h *Harness) RunAccess(ctx context.Context, urls []string, fetchInterval time.Duration) {
	// When the testharness uses an individual thing's DataURL
	// the resulting slice of things **MUST** contain a Thing with that same individual DataURL.
	// It **MAY** also contain other things if that single fetch happens to also retrieve them.

	fmt.Printf("########### Checking Access for: %s ########### \n", h.indexer.UID())

	fmt.Printf("CHECKING FOR ROBOTS.TXT FOR ALL URLS\n")
	allAllowed, allowErr := checkURLs(urls)
	if allowErr != nil {
		fmt.Printf("there is error from checkURLs %s\n", allowErr)
		return
	}

	if !allAllowed {
		fmt.Printf("the URLs are blocked by robots.txt\n")
		return
	}

	timeout := time.Duration(60) * time.Second
	clientFetch := thingfulx.NewClient("thingful", timeout)
	timeProvider := thingfulx.NewMockTimeProvider(time.Now())
	foundUniqueURL := 0 // to store how many have we found
	for _, u := range urls {

		foundUniqueURL = 0

		fmt.Printf("Fetching:  %s\n", u)
		bytes, err := h.indexer.Fetch(ctx, u, clientFetch)

		if err != nil {
			fmt.Printf("## ERROR from Fetch: %s\n", err.Error())
			fetchError = append(fetchError, err.Error())
			fetchErrorCount++
		}

		things, err := h.indexer.Parse(bytes, u, timeProvider)

		if err != nil {
			fmt.Printf("## ERROR from Parse: %s\n", err.Error())
			fetchError = append(fetchError, err.Error())
			fetchErrorCount++
		} else {

			if len(things) == 0 {
				fmt.Printf("ERROR this URL: %s returns nothing\n", u)
			} else {
				for i := 0; i < len(things); i++ {
					if u == things[i].Endpoint.URL { // check if one of "things" contain the same urls that used to access it
						foundUniqueURL++
						fmt.Printf("found same unique URL: %s \n", things[i].Endpoint.URL)
					}
				}
			}

		}

		if foundUniqueURL == 1 { // we only expect to find 1 copy of dataURL
			fmt.Printf("SUCCESS found one match: \n")
		} else if foundUniqueURL > 1 { // if we find more than 1, give warning
			fmt.Printf("ERROR found %d of the same unique URL: \n", foundUniqueURL)
		} else if foundUniqueURL == 0 { // also if we can't find any, give warning
			fmt.Printf("ERROR can't find anything that match: %s \n", u)
		}

	}

}

func min(number []int) int {

	min := 300

	for _, n := range number {
		min = int(math.Min(float64(min), float64(n)))
	}

	return min

}

func checkURLs(urls []string) (bool, error) {

	allAllowed := true
	robotsAddress := ""
	robots, err := robotstxt.FromString("User-agent: *\nDisallow:") //init robots
	if err != nil {
		return false, err
	}

	if whiteListed {

		fmt.Println("THIS PROVIDER IS WHITELISTED, IGNORING ROBOTS.TXT CHECK")
		allAllowed = true

	} else {

		for _, u1 := range urls {

			u, err := url.Parse(u1)
			if err != nil {
				return false, err
			}

			newRobotsAddress := u.Scheme + "://" + u.Host + "/robots.txt" // robots.txt address

			if newRobotsAddress != robotsAddress { // check if this is the same robots.txt, if not request new one

				robotsAddress = newRobotsAddress
				resp, err := http.Get(robotsAddress)
				if err != nil {
					return false, err
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				robots, err = robotstxt.FromBytes(body)

			}

			allow := robots.TestAgent(u.Path, "thingful") // then check if this path is allowed or not
			if !allow {
				fmt.Printf("%s is NOT allowed\n", u1)
				allAllowed = false
			}

		}

	}

	return allAllowed, nil

}
