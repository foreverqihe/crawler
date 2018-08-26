package core

import (
	"crypto/tls"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestRetrieveByDepth(t *testing.T) {

	t.Log("Testing retrieve")

	defer gock.Off()

	client := createClient()

	h := Handler{client}

	createDepthMocks()
	head := h.Retrieve("http://foo.com", 1)
	st.Refute(t, head, nil)
	st.Expect(t, head.URL, "http://foo.com")
	st.Expect(t, head.Title, "foo title")
	st.Expect(t, len(head.Nodes), 1)
	barNode := head.Nodes[0]
	st.Expect(t, barNode.URL, "http://bar.com")
	st.Expect(t, barNode.Depth, 1)
	st.Expect(t, barNode.Title, "") // it should not retrieve bar.com
	st.Expect(t, len(barNode.Nodes), 0)

	createDepthMocks() // mock again
	head = h.Retrieve("http://foo.com", 2)
	st.Refute(t, head, nil)
	st.Expect(t, head.URL, "http://foo.com")
	st.Expect(t, head.Title, "foo title")
	st.Expect(t, len(head.Nodes), 1)
	barNode = head.Nodes[0]
	st.Expect(t, barNode.URL, "http://bar.com")
	st.Expect(t, barNode.Depth, 1)
	st.Expect(t, barNode.Title, "bar title")
	st.Expect(t, len(barNode.Nodes), 1)

	barAgainNode := barNode.Nodes[0]
	st.Expect(t, barAgainNode.URL, "http://baragain.com")
	st.Expect(t, barAgainNode.Depth, 2)
	st.Expect(t, len(barAgainNode.Nodes), 0) // it shouldn't retrieve baragain.com
	st.Expect(t, barAgainNode.Title, "")

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
