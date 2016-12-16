# testharness

Light weight tool for fetcher develpers to run real tests locally without Pomelo

### How to implement with **existing** Fetcher
* update `glide.yaml` file to change thingfulx version to 0.1.9 and add testharness
```
- package: github.com/thingful/thingfulx
  version: 0.1.9
- package: github.com/thingful/testharness
```
* pull latest `bitbucket.org/thingful/template` at the moment it's `add-harness` branch
* copy `cmd` folder from `bitbucket.org/thingful/template` to the fetcher folder that you are testing
* replace `makefile` and `.gitignore` with the new ones from `bitbucket.org/thingful/template`
* change the word `template` in the main.go to the package name that you are testing
    * one in import `"bitbucket.com/thingful/template"`
    * another one in main function `harness, err := testharness.Register(template.NewFetcher)`
* run `glide up`
* run by `make harness`
* if API_KEY is needed do `API_KEY="123456" make harness`


### How to implement with **new** Fetcher
* make sure that new fetcher is based on latest `template`
* change the `template` in the main.go to the package name that you are testing
    * one in import `"bitbucket.com/thingful/template"`
    * another one in main function `harness, err := testharness.Register(template.NewFetcher)`
* run `glide up`
* run by `make harness`

## How to use
The code that runs TestHarness is at `com/harness/main.go`. This code is called by `make harness` command

### Register
Init by passing the fetcher and whitelisted bool to Register method.

`Register(builder thingfulx.FetcherBuilder, whitelisted bool)`

If this fetcher's whitelisted = `true`, it means that testHarness will ignore robots.txt completely. It is recommended to set this to `false` unless you are sure that we are allowed to ignore `robots.txt`

For example: 
`harness, err := testharness.Register(fetcherName.NewFetcher, false)`


There are 3 main features, these functions can be used separately.

### RunAll
RunAll runs URL method of this fetcher first

Then use the result URLs to fetch

Then for each URL to fetch, get one dataURL from first thing returned to access
```
fetchAllInterval := time.Duration(3) * time.Second // interval between each Fetch
urlsToFetch := 10                                  // total of urls to Fetch
harness.RunAll(context.Background(), fetchAllInterval, urlsToFetch)
```


### RunFetch
RunFetch only fetch the specified URLs
```
urls := []string{ // this is the URL produced by Fetcher's URL method
  "http://marlin.casa.ucl.ac.uk/~ollie/bikesapi/load.php?scheme=london",
}
fetchInterval := time.Duration(5) * time.Second // interval between each Fetch
harness.RunFetch(context.Background(), urls, fetchInterval)
```

### RunAccess
RunAccess fetch the dataURLs, then check that ONLY one thing returned from Fetch has the same dataURL
```
dataUrls := []string{ // this is the unique dataURL of each thing to fetch
  "http://marlin.casa.ucl.ac.uk/~ollie/bikesapi/load.php?scheme=london#id=73",
}
fetchInterval := time.Duration(5) * time.Second // interval between each Access
harness.RunAccess(context.Background(), dataUrls, fetchInterval)
```

