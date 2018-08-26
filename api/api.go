package api

import (
	"encoding/json"
	"net/http"

	"github.com/foreverqihe/crawler/core"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// swagger:operation POST /v1/crawl crawl
//
// crawl the url with a limit of depth
//
// ---
// produces:
// - application/json
// parameters:
// - name: request_json
//   in: body
//   description: the request url and depth
//   required: true
//   schema:
//     type: object
//     properties:
//       url:
//         type: string
//         description: the root url to retrieve
//       depth:
//         type: integer
//         description: 0 doesn't retrieve anything, 1 retrieves the root url only, 2 retries level one leaves, etc.
// responses:
//  '200':
//    description: successful and body will contain a json tree
//  '400':
//    description: bad request, the request json maybe malformatted
//  '500':
//    description: server internal error
func (c *CrawlerAPI) craw(w http.ResponseWriter, r *http.Request) {
	var req = struct {
		URL   string `json:"url"`
		Depth int    `json:"depth"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Errorln("cannot parse body")
		http.Error(w, "cannot parse body json", http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{"url": req.URL, "depth": req.Depth}).Infoln("request received")

	root := c.Handler.Retrieve(req.URL, req.Depth)
	jsonBytes, err := json.MarshalIndent(root, "", " ")
	if err != nil {
		http.Error(w, "failed to stringify json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Crawler
type CrawlerAPI struct {
	Handler *core.Handler
}

// Start point
func (s *CrawlerAPI) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/v1/crawl", s.craw).Methods("POST")

	corsOriginObj := handlers.AllowedOrigins([]string{"*"})
	corsMethodsObj := handlers.AllowedMethods([]string{"POST", "GET", "OPTION"})
	corsAllowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type"})
	server := &http.Server{
		Handler: handlers.CORS(corsOriginObj, corsMethodsObj, corsAllowedHeaders)(router),
		Addr:    ":8080",
	}
	return server.ListenAndServe()
}
