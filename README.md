# testharness

Light weight tool for fetcher develpers to run real tests locally without Pomelo

### How to use with **existing** Fetcher
* update `glide.yaml` file to change thingfulx version to 0.1.8 and add testharness
```
- package: github.com/thingful/thingfulx
  version: 0.1.8
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


### How to use with **new** Fetcher
* make sure that new fetcher is based on latest `template`
* change the `template` in the main.go to the package name that you are testing
    * one in import `"bitbucket.com/thingful/template"`
    * another one in main function `harness, err := testharness.Register(template.NewFetcher)`
* run `glide up`
* run by `make harness`

