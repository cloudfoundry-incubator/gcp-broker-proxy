package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/urfave/negroni"

	"code.cloudfoundry.org/gcp-broker-proxy/auth"
	"code.cloudfoundry.org/gcp-broker-proxy/oauth"
	"code.cloudfoundry.org/gcp-broker-proxy/proxy"
	"code.cloudfoundry.org/gcp-broker-proxy/startupchecker"
	"code.cloudfoundry.org/gcp-broker-proxy/token"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	username, password, brokerURLString, serviceAccountJSON := getRequiredEnvs()

	brokerURL, err := url.ParseRequestURI(brokerURLString)
	if err != nil {
		log.Fatal(fmt.Sprintf("BROKER_URL must be a valid URL: %s", brokerURLString))
	}

	tokenFetcher, err := oauth.NewGCPOAuth(serviceAccountJSON)
	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid SERVICE_ACCOUNT_JSON: %s", err))
	}

	client := http.Client{}

	startupChecker := startupchecker.NewChecker(brokerURL, tokenFetcher, &client)

	err = startupChecker.Perform()
	if err != nil {
		log.Fatal("Failed startup checks: " + err.Error())
	}
	fmt.Println("Startup checks passed")

	basicAuth := auth.BasicAuth(username, password)
	reverseProxy := proxy.ReverseProxy(brokerURL)
	tokenHandler := token.TokenHandler(tokenFetcher)

	n := negroni.New()

	logger := negroni.NewLogger()
	logger.SetFormat("{{.Status}} | {{.Method}} {{.Path}} {{.Request.URL.RawQuery}} | \t {{.Duration}} \n")

	n.Use(logger)
	n.Use(basicAuth)
	n.Use(tokenHandler)
	n.Use(reverseProxy)

	fmt.Printf("About to listen on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, n))
}

func getRequiredEnvs() (username, password, brokerURL, serviceAccountJSON string) {
	var missingEnvs []string

	getRequiredEnv := func(env string) string {
		parsedEnv := os.Getenv(env)
		if parsedEnv == "" {
			missingEnvs = append(missingEnvs, env)
		}
		return parsedEnv
	}

	username = getRequiredEnv("USERNAME")
	password = getRequiredEnv("PASSWORD")
	brokerURL = getRequiredEnv("BROKER_URL")
	serviceAccountJSON = getRequiredEnv("SERVICE_ACCOUNT_JSON")

	if len(missingEnvs) != 0 {
		errMsg := fmt.Sprintf("Missing %s environment variable(s)", strings.Join(missingEnvs, ", "))
		log.Fatal(errMsg)
	}

	return
}
