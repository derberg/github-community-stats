# GitHub Community Stats
Golang based learning project that fetches some basic stats about community of a given GitHub organization

## Development

### Install dependencies

This project uses `dep` as a dependency manager. To install all required dependencies, use the following command:
```bash
dep ensure -vendor-only
```

### Run project

GITHUB_ORG_ID=MDEyOk9yZ2FuaXphdGlvbjM5MTUzNTIz GITHUB_ORG=kyma-project GITHUB_REPO=kyma GITHUB_TOKEN={YOUR_TOKEN} gorun main.go

## Current state of the project

It is in progress and for now hardcoded to get data for only one specific organization.