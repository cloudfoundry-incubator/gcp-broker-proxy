# Google Cloud Platform Proxy Service Broker
[![Build Status](https://travis-ci.org/cloudfoundry-incubator/gcp-broker-proxy.svg?branch=master)](https://travis-ci.org/cloudfoundry-incubator/gcp-broker-proxy)

**Note**: This repository should be imported as code.cloudfoundry.org/gcp-broker-proxy.


This broker proxies requests to Google's hosted service broker. It handles the OAuth flow and allows the GCP
to be registered in CloudFoundry.

### Installation
```
go get -u code.cloudfoundry.org/gcp-broker-proxy
```

### Deploying to Cloud Foundry
1. Configure the broker by setting the environment variables in the `manifest.yml`.
1. `make build-linux`
1. `cf push`

### Development

#### Test
```
export TEST_GCP_SERVICE_ACCOUNT_JSON=$(cat service_account_json.json)
make test
```

#### Build
```
make build
```
