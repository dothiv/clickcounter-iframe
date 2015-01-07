package clickcounteriframe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	hivdomainstatus "github.com/dothiv/hiv-domain-status"
)

type AdminController struct {
	domainRepo DomainRepositoryInterface
}

func NewAdminController(d DomainRepositoryInterface) (c *AdminController) {
	c = new(AdminController)
	c.domainRepo = d
	return
}

func (c *AdminController) DomainHandler(w http.ResponseWriter, r *http.Request, matches []string) {
	w.Header().Add("X-Click-Counter-Iframe-Version", VERSION)
	if r.Method == "DELETE" {
		c.deleteDomain(w, r, matches)
		return
	}
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		hivdomainstatus.HttpProblem(w, http.StatusBadRequest, "Expected application/json")
		return
	}

	b, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		hivdomainstatus.HttpProblem(w, http.StatusInternalServerError, fmt.Sprintf("failed to read request body:", bodyErr.Error()))
		return
	}
	var data struct {
		Redirect    string
		LandingPage *LandingPage `json:"landingPage"`
	}
	unmarshalErr := json.Unmarshal(b, &data)
	if unmarshalErr != nil {
		hivdomainstatus.HttpProblem(w, http.StatusBadRequest, fmt.Sprintf("failed to read json: %s", unmarshalErr.Error()))
		return
	}

	redirect, urlErr := url.Parse(data.Redirect)
	if urlErr != nil {
		hivdomainstatus.HttpProblem(w, http.StatusBadRequest, fmt.Sprintf("Invalid redirect url provided: %s", data.Redirect))
		return
	}

	domain, domainErr := c.domainRepo.FindByName(matches[1])
	created := false
	if domainErr != nil {
		domain = new(Domain)
		domain.Name = matches[1]
		created = true
	}
	if len(redirect.String()) > 0 {
		domain.Redirect.String = redirect.String()
		domain.Redirect.Valid = true
	}

	domain.LandingPage = data.LandingPage
	c.domainRepo.Persist(domain)
	if created {
		w.Header().Add("Location", r.URL.String())
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *AdminController) deleteDomain(w http.ResponseWriter, r *http.Request, matches []string) {
	domain, domainErr := c.domainRepo.FindByName(matches[1])
	if domainErr != nil {
		hivdomainstatus.HttpProblem(w, http.StatusNotFound, fmt.Sprintf("domain not found: %s", matches[1]))
		return
	}
	deleteErr := c.domainRepo.Remove(domain)
	if deleteErr != nil {
		hivdomainstatus.HttpProblem(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete domain: %s! %s", matches[1], deleteErr.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
