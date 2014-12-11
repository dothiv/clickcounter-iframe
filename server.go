package clickcounteriframe

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func Serve(c *Config) (err error) {
	// Open DB
	db, err := sql.Open("postgres", c.DSN())
	if err != nil {
		return
	}

	log.Println(fmt.Sprintf("Starting server on port %d ...", c.Server.Port))

	domainRepo := NewDomainRepository(db)

	iframeCntrl := NewIframeController(domainRepo, c.Server.Hostname)
	adminCntrl := NewAdminController(domainRepo, c.Auth.Token)

	reHandler := new(RegexpHandler)
	reHandler.AddRoute("^/domain/([^/]+)$", adminCntrl.DomainHandler)
	reHandler.AddRoute("^/$", iframeCntrl.IframeHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", c.Server.Port), reHandler))
	return
}
