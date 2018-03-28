package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strconv"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	_ "code.cloudfoundry.org/gcp-broker-proxy"
)

var _ = Describe("GCP Broker Proxy", func() {
	var (
		session *gexec.Session
		port    string
		envs    []string
	)

	var brokerServer *httptest.Server
	var brokerURL string

	BeforeEach(func() {
		brokerServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "I'm a broker")
		}))

		brokerURL = brokerServer.URL
	})

	JustBeforeEach(func() {
		var err error

		cmd := exec.Command(gcpBrokerProxyBinary)
		cmd.Env = envs
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		session.Kill()
	})

	Describe("when the server is correctly configured", func() {
		BeforeEach(func() {
			port = strconv.Itoa(8081 + config.GinkgoConfig.ParallelNode)
			envs = []string{"PORT=" + port, "SERVICE_ACCOUNT_JSON=" + testServiceAccountJSON, "BROKER_URL=" + brokerURL, "USERNAME=admin", "PASSWORD=foo"}
		})

		It("logs that the server is about to start on a specific port", func() {
			Eventually(session).Should(Say("About to listen on port " + port))
		})

		It("does not exit", func() {
			Consistently(session).ShouldNot(gexec.Exit())
		})

		Context("when no port is specified", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON=" + testServiceAccountJSON, "BROKER_URL=" + brokerURL, "USERNAME=admin", "PASSWORD=foo"}
			})

			It("it starts on the default port of 8080", func() {
				Eventually(session).Should(Say("About to listen on port 8080"))
			})
		})
	})

	Describe("when the server is not correctly configured", func() {
		Context("when the server has not been provided service account information", func() {
			BeforeEach(func() {
				envs = []string{"BROKER_URL=" + brokerURL, "USERNAME=admin", "PASSWORD=foo"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the SERVICE_ACCOUNT_JSON param", func() {
				Eventually(session.Err).Should(Say("Missing SERVICE_ACCOUNT_JSON environment variable"))
			})
		})

		Context("when the server has not been provided broker url", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON={\"type\": \"service_account\"}", "USERNAME=admin", "PASSWORD=foo"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the BROKER_URL param", func() {
				Eventually(session.Err).Should(Say("Missing BROKER_URL environment variable"))
			})
		})

		Context("when the broker url is invalid", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON={\"type\": \"service_account\"}", "BROKER_URL=notaurl", "USERNAME=admin", "PASSWORD=foo"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the BROKER_URL param", func() {
				Eventually(session.Err).Should(Say("BROKER_URL must be a valid URL: notaurl"))
			})
		})

		Context("when the server has not been provided username", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON={\"type\": \"service_account\"}", "BROKER_URL=" + brokerURL, "PASSWORD=foo"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the USERNAME param", func() {
				Eventually(session.Err).Should(Say("Missing USERNAME environment variable"))
			})
		})

		Context("when the server has not been provided password", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON={\"type\": \"service_account\"}", "BROKER_URL=" + brokerURL, "USERNAME=admin"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the PASSWORD param", func() {
				Eventually(session.Err).Should(Say("Missing PASSWORD environment variable"))
			})
		})

		Context("when there are multiple missing parameters", func() {
			BeforeEach(func() {
				envs = []string{}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that it requires all missing params", func() {
				Eventually(session.Err).Should(Say("Missing USERNAME, PASSWORD, BROKER_URL, SERVICE_ACCOUNT_JSON environment variable\\(s\\)"))
			})
		})
	})
})
