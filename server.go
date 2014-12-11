package clickcounteriframe

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type route struct {
	re      *regexp.Regexp
	handler func(http.ResponseWriter, *http.Request, []string)
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) AddRoute(re string, handler func(http.ResponseWriter, *http.Request, []string)) {
	r := &route{regexp.MustCompile(re), handler}
	h.routes = append(h.routes, r)
}

func (h *RegexpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		matches := route.re.FindStringSubmatch(r.URL.Path)
		if matches != nil {
			route.handler(rw, r, matches)
			break
		}
	}
}

func Serve(c *Config) (err error) {
	// Open DB
	db, err := sql.Open("postgres", c.DSN())
	if err != nil {
		return
	}

	log.Println(fmt.Sprintf("Starting server on port %d ...", c.Server.Port))

	domainRepo := NewDomainRepository(db)

	iframeCntrl := NewIframeController(domainRepo, c.Server.Hostname)
	adminCntrl := NewAdminController(domainRepo)

	reHandler := new(RegexpHandler)
	reHandler.AddRoute("^/domain/([^/]+)$", adminCntrl.DomainHandler)
	reHandler.AddRoute("^/$", iframeCntrl.IframeHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", c.Server.Port), reHandler))
	return
}
