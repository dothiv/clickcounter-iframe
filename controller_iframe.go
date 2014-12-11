package clickcounteriframe

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type IframeController struct {
	domainRepo     DomainRepositoryInterface
	hostname       string
	cacheLifetTime time.Duration
}

func NewIframeController(d DomainRepositoryInterface, hostname string) (c *IframeController) {
	c = new(IframeController)
	c.domainRepo = d
	c.hostname = hostname
	c.cacheLifetTime = time.Minute * 30
	return
}

var iframeTpl = `<!DOCTYPE html>
<!--


            _|_|          _|_|
            _|_|          _|_|
            _|_|
            _|_|
            _|_|_|_|      _|_|  _|_|    _|_|
            _|_|_|_|_|    _|_|  _|_|    _|_|
            _|_|    _|_|  _|_|  _|_|    _|_|
            _|_|    _|_|  _|_|  _|_|    _|_|
            _|_|    _|_|  _|_|  _|_|    _|_|
      _|_|  _|_|    _|_|  _|_|    _|_|_|_|
      _|_|  _|_|    _|_|  _|_|      _|_|

      .hiv domains – The digital Red Ribbon

                  click4life.hiv

-->
<html>
<head>
    <title>{{.Name}}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style type="text/css">
        html, body {
            margin: 0;
            padding: 0;
            width: 100%;
            height: 100%;
            overflow: hidden;
        }

        #clickcounter-target-iframe {
            border: 0;
            width: 100%;
            height: 100%;
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
<iframe src="{{.Redirect}}" width="100%" height="100%" id="clickcounter-target-iframe"></iframe>
<script src="//dothiv-registry.appspot.com/static/clickcounter.min.js" type="text/javascript"></script>
</body>
</html>`

func (c *IframeController) IframeHandler(w http.ResponseWriter, r *http.Request, matches []string) {
	w.Header().Add("X-Click-Counter-Iframe-Version", VERSION)
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	hivDomainName := c.getHivDomainName(r) + ".hiv"
	domain, err := c.domainRepo.FindByName(hivDomainName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Domain not found:", hivDomainName)
		return
	}

	w.Header().Add("Cache-Control", fmt.Sprintf("public, s-maxage=%d", int64(c.cacheLifetTime/time.Second)))
	w.Header().Add("Expires", time.Now().Add(c.cacheLifetTime).Format(http.TimeFormat))
	w.Header().Add("Last-Modified", domain.Updated.Format(http.TimeFormat))

	t := template.Must(template.New("iframe").Parse(iframeTpl))
	err = t.Execute(w, domain)
	if err != nil {
		log.Println("failed to parse template:", err)
	}
}

func (c *IframeController) getHivDomainName(r *http.Request) (hostname string) {
	var hostNameMatch = regexp.MustCompile(`([^\.]+)\.` + c.hostname)
	hostname = hostNameMatch.FindStringSubmatch(strings.Split(r.Host, ":")[0])[1]
	return
}
