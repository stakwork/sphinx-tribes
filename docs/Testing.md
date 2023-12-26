
# Testing on sphinx-tribes

There are two main tests that we run in this project
- Golang/ our backend tests
- Jest unit and component tests/ frontend test

We have git actions for both of these which can be found in `./.github/workflows/prjob_tests.yml`

# Golang/ our backend tests
#### Run the tests
Run this command you may need to install `cover` first run
- `go get golang.org/x/tools/cmd/cover`

To run the tests run
- `go test ./... -tags mock -race -v`
#### Run with coverage
- `go test ./... -tags mock -race -v -coverprofile=coverage.out && ./cover-check.sh coverage.out <min coverage amount>`
#### Other details
- We have a github pr job in `./.github/workflows/prjob_tests.yml` this includes the jest test runner.
- Our tests can be identified in the codebase if you see a file ending in `<filename>_test.go`


## Jest unit and component tests/ frontend tests
In our frontend testing suite we have two types of tests that we run 
- Unit tests for files that end in `.ts`
- Component tests for files that end in `.tsx`

You can find configuration in `./frontend/app/jest.config.json`
this includes which files are include and our minimum test coverage amounts.
#### Run test (this is with coverage)
- navigate to `./frontend/app`
- then type `yarn run test`

