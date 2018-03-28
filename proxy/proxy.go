package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

//go:generate counterfeiter . TokenRetriever
type TokenRetriever interface {
	GetToken() (*oauth2.Token, error)
}

//go:generate counterfeiter . HTTPDoer
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Proxy struct {
	brokerURL      string
	tokenRetriever TokenRetriever
	httpDoer       HTTPDoer
}

func NewProxy(brokerURL string, tr TokenRetriever, httpDoer HTTPDoer) Proxy {
	return Proxy{
		brokerURL:      brokerURL,
		tokenRetriever: tr,
		httpDoer:       httpDoer,
	}
}

func (p *Proxy) PerformStartupChecks() error {
	token, err := p.tokenRetriever.GetToken()
	if err != nil {
		return errors.Wrap(err, "Failed obtaining oauth token")
	}

	req, err := http.NewRequest("GET", p.brokerURL+"/v2/catalog", nil)

	if err != nil {
		return errors.Wrap(err, "Failed to create request")
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("x-broker-api-version", "2.14")

	res, err := p.httpDoer.Do(req)

	if err != nil {
		return errors.Wrap(err, "Failed to make request to the broker")
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		bodyBytes, err := ioutil.ReadAll(res.Body)
		var bodyString string
		if err != nil {
			bodyString = "Could not read body"
		} else {
			bodyString = string(bodyBytes)
		}
		return fmt.Errorf("Broker did not respond successfully. status: %d body: %s", res.StatusCode, bodyString)
	}

	return err
}
