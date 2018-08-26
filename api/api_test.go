package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/foreverqihe/crawler/core"
	"github.com/gorilla/mux"
	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestMain(m *testing.M) {
	before()
	code := m.Run()
	after()
	os.Exit(code)
}

func TestRetrieveDepthZero(t *testing.T) {
	t.Log("it should not retrieve the root when depth is 0")
	{
		reqObj := struct {
			URL   string `json:"url"`
			Depth int    `json:"depth"`
		}{
			URL:   "http://foo.com",
			Depth: 0,
		}
		reqBody, err := json.Marshal(reqObj)

		st.Assert(t, err, nil)
		req, err := http.NewRequest("POST", "/v1/crawl", bytes.NewReader(reqBody))
		st.Assert(t, err, nil)

		var respObj core.Node
		rw := httptest.NewRecorder()
		router.ServeHTTP(rw, req)
		st.Assert(t, rw.Code, http.StatusOK)
		err = json.NewDecoder(rw.Body).Decode(&respObj)
		st.Assert(t, err, nil)
		st.Expect(t, respObj.Title, "")
	}
}

func TestRetrieveDepth(t *testing.T) {
	t.Log("it should retrieve by depth 1")
	{
		createDepthMocks()
		reqObj := struct {
			URL   string `json:"url"`
			Depth int    `json:"depth"`
		}{
			URL:   "http://foo.com",
			Depth: 1,
		}
		reqBody, err := json.Marshal(reqObj)

		st.Assert(t, err, nil)
		req, err := http.NewRequest("POST", "/v1/crawl", bytes.NewReader(reqBody))
		st.Assert(t, err, nil)

		var respObj core.Node
		rw := httptest.NewRecorder()
		router.ServeHTTP(rw, req)
		st.Assert(t, rw.Code, http.StatusOK)
		err = json.NewDecoder(rw.Body).Decode(&respObj)
		st.Assert(t, err, nil)
		st.Expect(t, respObj.Title, "foo title")
		st.Expect(t, len(respObj.Nodes), 1)
		barNode := respObj.Nodes[0]
		st.Expect(t, barNode.URL, "http://bar.com")
		st.Expect(t, barNode.Title, "")
		st.Expect(t, len(barNode.Nodes), 0)
	}

	t.Log("it should retrieve by depth 2")
	{
		createDepthMocks()
		reqObj := struct {
			URL   string `json:"url"`
			Depth int    `json:"depth"`
		}{
			URL:   "http://foo.com",
			Depth: 2,
		}
		reqBody, err := json.Marshal(reqObj)

		st.Assert(t, err, nil)
		req, err := http.NewRequest("POST", "/v1/crawl", bytes.NewReader(reqBody))
		st.Assert(t, err, nil)

		var respObj core.Node
		rw := httptest.NewRecorder()
		router.ServeHTTP(rw, req)
		st.Assert(t, rw.Code, http.StatusOK)
		err = json.NewDecoder(rw.Body).Decode(&respObj)
		st.Assert(t, err, nil)
		st.Expect(t, respObj.Title, "foo title")
		st.Expect(t, len(respObj.Nodes), 1)
		barNode := respObj.Nodes[0]
		st.Expect(t, barNode.URL, "http://bar.com")
		st.Expect(t, barNode.Title, "bar title")
		st.Expect(t, len(barNode.Nodes), 1)
		barAgainNode := barNode.Nodes[0]
		st.Expect(t, barAgainNode.URL, "http://baragain.com")
		st.Expect(t, barAgainNode.Title, "")
		st.Expect(t, len(barAgainNode.Nodes), 0)
	}
}

var router *mux.Router

func before() {

	router = mux.NewRouter()
	handler := core.NewHander()
	gock.InterceptClient(handler.Client)
	api := CrawlerAPI{
		Handler: handler,
	}
	router.HandleFunc("/v1/crawl", api.craw).Methods("POST")
}

func after() {

}
func createDepthMocks() {
	foo := `
<html>
  <head>
    <title>foo title</title>
  </head>
  <body>
    <p>
     Test on this dummy link: <a "href"="http://bar.com">Bar</a>
    </p>
  </body>
</html>
`
	bar := `
<html>
  <head>
    <title>bar title</title>
  </head>
  <body>
    <p>
    <a 'href'="http://baragain.com">Bargain</a> <strong>it supports both ' and "
    </p>
</html>
`
	gock.New("http://foo.com").
		Get("/").
		Reply(200).BodyString(foo)
	gock.New("http://bar.com").
		Get("/").
		Reply(200).
		BodyString(bar)
	gock.New("http://baragain.com").Get("/").Reply(200).BodyString("nothing")
}
func createClient() *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 20 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 20 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second, // default client never timeout
	}
	gock.InterceptClient(client)
	return client
}
