package clickcounteriframe

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type IframeController struct {
	domainRepo     DomainRepositoryInterface
	cacheLifetTime time.Duration
}

func NewIframeController(d DomainRepositoryInterface) (c *IframeController) {
	c = new(IframeController)
	c.domainRepo = d
	c.cacheLifetTime = time.Minute * 30
	return
}

func (c *IframeController) IframeHandler(w http.ResponseWriter, r *http.Request, matches []string) {
	w.Header().Add("X-Click-Counter-Iframe-Version", VERSION)
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	secondLevelName, err := c.getSecondLevelName(r)
	hivDomainName := secondLevelName + ".hiv"
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err.Error())
		return
	}
	domain, err := c.domainRepo.FindByName(hivDomainName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Domain not found:", hivDomainName)
		return
	}

	w.Header().Add("Cache-Control", fmt.Sprintf("public, s-maxage=%d", int64(c.cacheLifetTime/time.Second)))
	w.Header().Add("Expires", time.Now().Add(c.cacheLifetTime).Format(http.TimeFormat))
	w.Header().Add("Last-Modified", domain.Updated.Format(http.TimeFormat))

	var tpl *template.Template
	if len(domain.Redirect) > 0 {
		tpl, err = getIframeTemplate()
	} else {
		tpl, err = getLandingPageTemplate()
		if &domain.LandingPage.Tweet != nil {
			domain.LandingPage.TweetEncoded = url.QueryEscape(domain.LandingPage.Tweet)
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to load template")
		log.Println(err.Error())
		return
	}
	err = tpl.Execute(w, domain)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to parse template")
		log.Println(err.Error())
		return
	}
}

var iframeTpl *template.Template

func getIframeTemplate() (tpl *template.Template, err error) {
	if iframeTpl == nil {
		iframeTpl, err = loadTemplate("iframe", "./templates/iframe.html")
		if err != nil {
			return
		}
	}
	return iframeTpl, err
}

var landingPageTpl *template.Template

func getLandingPageTemplate() (tpl *template.Template, err error) {
	//if landingPageTpl == nil {
	landingPageTpl, err = loadTemplate("landingpage", "./templates/landingpage.html")
	if err != nil {
		return
	}
	//}
	return landingPageTpl, err
}

func loadTemplate(ident string, filename string) (tpl *template.Template, err error) {
	tplSource, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	tpl = template.Must(template.New(ident).Parse(string(tplSource)))
	return
}

func (c *IframeController) getSecondLevelName(r *http.Request) (secondLevelName string, err error) {
	var hostNameMatch = regexp.MustCompile(`([^\.]+)\.hiv$`)
	domainName := strings.Split(r.Host, ":")[0]
	match := hostNameMatch.FindStringSubmatch(domainName)
	if len(match) == 0 {
		err = fmt.Errorf("Not a .hiv domain name: %s", domainName)
		return
	}
	secondLevelName = match[1]
	return
}
