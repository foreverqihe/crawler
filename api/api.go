package api

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type CrawlerApi struct {
}

func (s *CrawlerApi) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/v1/crawl", s.craw).Methods("GET").Queries("url", "{url}")

	corsOriginObj := handlers.AllowedOrigins([]string{"*"})
	corsMethodsObj := handlers.AllowedMethods([]string{"GET"})
	corsAllowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type"})
	server := &http.Server{
		Handler: handlers.CORS(corsOriginObj, corsMethodsObj, corsAllowedHeaders)(router),
		Addr:    ":8080",
	}
	return server.ListenAndServe()
}

func (c *CrawlerApi) craw(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	log.WithFields(log.Fields{"url": url}).Infoln("request received")

}
