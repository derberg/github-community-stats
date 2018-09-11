# GitHub Community Stats
Golang based learning project that fetches some basic stats about community of a given GitHub organization

## Development

### Install dependencies

This project uses `dep` as a dependency manager. To install all required dependencies, use the following command:
```bash
dep ensure -vendor-only
```

### Run project

GITHUB_GRAPHQL_TEST_TOKEN="{YOUR_SECRET_TOKEN}" go run main.go

## Current state of the project

It is in progress and for now hardcoded to get data for only one specific organization.