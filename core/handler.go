package core

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type Handler struct {
	Client *http.Client
}

func NewHander() *Handler {
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
	return &Handler{client}
}

type Node struct {
	URL   string  `json:"url"`
	Title string  `json:"title"`
	Nodes []*Node `json:"nodes"`
	Depth int
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func (h *Handler) Retrieve(rootUrl string, depth int) *Node {

	var n int
	tasks := make(chan []*Node)

	unseen := make(chan *Node)

	head := &Node{
		URL:   rootUrl,
		Depth: 0,
	}
	if depth <= 0 {
		log.WithField("depth", depth).Debugln("depth <= 0, do not retrieve")
		return head
	}
	n++
	go func() {
		tasks <- []*Node{head}
	}()

	for i := 0; i < (5 * runtime.NumCPU()); i++ {
		go func() {
			for link := range unseen {
				h.crawl(link)
				go func(link *Node) { tasks <- link.Nodes }(link)
			}
		}()
	}

	seen := make(map[string]bool)

	for ; n > 0; n-- {
		list := <-tasks
		for _, node := range list {
			if !seen[node.URL] && node.Depth < depth {
				seen[node.URL] = true
				n++
				unseen <- node
			}
		}
	}

	return head
}

func (h *Handler) crawl(node *Node) {

	log.Debugln("crawling", node.URL)
	resp, err := h.Client.Get(node.URL)
	if err != nil {
		log.WithError(err).WithField("url", node.URL).Errorln("failed to get url")
		return
	}

	defer resp.Body.Close()

	links, title := h.getAllLinks(resp.Body, node.URL)
	log.Debugln("got title", title)
	log.Debugf("links: %+v", links)
	node.Title = title
	for _, link := range links {
		if strings.HasPrefix(link, "http") {
			// http only, do not support mailto:
			log.WithField("link", link).Debugln("got a link")
			newNode := &Node{
				URL:   link,
				Depth: node.Depth + 1,
			}
			node.Nodes = append(node.Nodes, newNode)
		}
	}
}

func (h *Handler) getAllLinks(body io.Reader, url string) (links []string, title string) {

	z := html.NewTokenizer(body)
	depth := 0
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.TextToken:
			if depth > 0 {
				title = string(z.Text())
			}
		case html.StartTagToken, html.EndTagToken:
			tn, hasAttr := z.TagName()
			if len(tn) == 5 && string(tn) == "title" {
				if tt == html.StartTagToken {
					depth++
				} else {
					depth--
				}
			} else if len(tn) == 1 && tn[0] == 'a' && hasAttr {

				if tt == html.StartTagToken {
					for {
						key, val, hasMore := z.TagAttr()
						log.Debugln("found a link, key =", string(key))
						if h.removeQuotes(string(key)) == "href" {

							links = append(links, h.resolveLink(string(val), url))
						}
						if !hasMore {
							break
						}
					}
				}
			}

		}
	}

}

func (h *Handler) removeQuotes(s string) string {
	if len(s) > 0 && (s[0] == '"' || s[0] == '\'') {
		s = s[1:]
	}
	if len(s) > 0 && (s[len(s)-1] == '"' || s[len(s)-1] == '\'') {
		s = s[:len(s)-1]
	}
	return s
}

// 1. drop # and anything after it
// 2. resolve to absolute url
func (h *Handler) resolveLink(l string, base string) string {
	if strings.Contains(l, "#") {
		var index int
		for n, str := range l {
			if strconv.QuoteRune(str) == "'#'" {
				index = n
				break
			}
		}
		l = l[:index]
	}
	uri, err := url.Parse(l)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	log.Debugln("resolved link", uri.String())
	return uri.String()
}
