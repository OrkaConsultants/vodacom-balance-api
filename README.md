# Vodacom Balances API

This code is not developed by Vodacom and it is only an interpretation of their API.

The intention of this GoLang API is to query balances of SIM cards through an easy REST API.

## Getting Started

1. Install golang dependency management tool: `go get -u github.com/golang/dep/cmd/dep`
2. Install dependencies `dep ensure -v`
3. Ensure configuration is correct in `cmd/vodacom-api/bootstrap.yml` and `cmd/vodacom-api/vodacom-api.yml`
4. Change directory to `cd cmd/vodacom-api/` Execute `go run main.go`

## Configuration

### `bootstrap.yml`

Startup Configuration

``` yml
# IGNORE THIS IF YOU DO NOT INTENT TO RUN THIS AS A MICROSERVICE
eureka: # Spring Boot Eureka
  url: https://localhost
  port: 8761
  enabled: false

app:
  name: vodacom-api
  port: 8080
  authentication: false # JWT Authentication for this application

# IGNORE THIS IF YOU DO NOT INTENT TO RUN THIS AS A MICROSERVICE
cloud-config: # Spring cloud condig
  enabled: false
  profile: demo
  service-name: config-service # Config service name
```

### `vodacom-api.yml`

Runtime Configuration

``` yml
jwt-config: # Fromat of JWT Authentication Ignore if it is disabled in bootstrap
  TOKEN_PREFIX: Bearer
  HEADER_STRING: Authorization

vodacom-api:
  username: your@email.com
  password: greatPassword1234
  login-uri: https://myvodacom.secure.vodacom.co.za/cloud/login/v3/login # DO NOT CHANGE
  balance-uri: https://www.vodacom.co.za/cloud/rest/balances/v2/bundleBalances/{number} # DO NOT CHANGE
```

## How does this work?

1. Run the application.
2. The application will start up on http://localhost:8080.
3. Your login details specified in `vodacom-api.yml` will be used to receive an API token to get the balances.
4. Open your favorite API Development program, like Postman.
5. Perform an HTTP GET request to http://localhost:8080/api/balance?number=278200LEGGEH (The number specified here must start with 27 (2782... and NOT 082) and it must be linked to your account)
6. A list of services will be returned in the following format:

``` JSON
{
    "showDial": true,
    "showBuyButton": true,
    "remaining": 2300.0,
    "total": 0.0,
    "preText": "Airtime",
    "midText": "R 23.00",
    "postText": "Remaining"
},
{
    "showDial": true,
    "showBuyButton": true,
    "remaining": 234470.0,
    "total": 256000.0,
    "preText": "Data",
    "midText": "228.97MB",
    "postText": "Remaining"
},
{
    "showDial": true,
    "showBuyButton": true,
    "remaining": 3600.0,
    "total": 0.0,
    "preText": "Voice",
    "midText": "0m0s",
    "postText": "Remaining"
},
{
    "showDial": true,
    "showBuyButton": true,
    "remaining": 0.0,
    "total": 0.0,
    "preText": "SMS",
    "midText": "0",
    "postText": "Remaining"
},
{
    "showDial": true,
    "showBuyButton": true,
    "remaining": 0.0,
    "total": 0.0,
    "preText": "MMS",
    "midText": "0",
    "postText": "Remaining"
}
```

## Issues

Visit the issues section on github.
