package clickcounteriframe

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupAdminTest(t *testing.T) (cntrl *AdminController, domainRepo DomainRepositoryInterface) {

	c, configErr := NewConfig("config.ini")
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())

	domainRepo = NewDomainRepository(db)
	cntrl = NewAdminController(domainRepo, "letmein")
	db.Exec("TRUNCATE domain RESTART IDENTITY")

	return
}

func TestThatItCreatesDomainConfig(t *testing.T) {
	assert := assert.New(t)

	cntrl, repo := SetupAdminTest(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntrl.DomainHandler(w, r, regexp.MustCompile("^/domain/([^/]+)$").FindStringSubmatch(r.URL.Path))
	}))
	defer ts.Close()

	var data = []byte(`{"redirect":"http://example.com/"}`)
	req, err := http.NewRequest("PUT", ts.URL+"/domain/example.hiv", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer letmein")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(http.StatusCreated, resp.StatusCode)

	// Verify
	d, findErr := repo.FindByName("example.hiv")
	assert.Nil(findErr)

	assert.Equal(1, d.Id)
	assert.Equal("example.hiv", d.Name)
	assert.Equal("http://example.com/", d.Redirect)
}

func TestThatItUpdatesDomainConfig(t *testing.T) {
	assert := assert.New(t)

	cntrl, repo := SetupAdminTest(t)
	d := new(Domain)
	d.Name = "acme.hiv"
	d.Redirect = "http://acme.info/"
	repo.Persist(d)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntrl.DomainHandler(w, r, regexp.MustCompile("^/domain/([^/]+)$").FindStringSubmatch(r.URL.Path))
	}))
	defer ts.Close()

	var data = []byte(`{"redirect":"http://acme.com/"}`)
	req, err := http.NewRequest("PUT", ts.URL+"/domain/acme.hiv", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer letmein")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify
	d, findErr := repo.FindByName("acme.hiv")
	assert.Nil(findErr)

	assert.Equal(1, d.Id)
	assert.Equal("acme.hiv", d.Name)
	assert.Equal("http://acme.com/", d.Redirect)
}

func TestThatItDeletesDomainConfig(t *testing.T) {
	assert := assert.New(t)

	cntrl, repo := SetupAdminTest(t)
	d := new(Domain)
	d.Name = "microsoft.hiv"
	d.Redirect = "http://microsoft.com/"
	repo.Persist(d)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntrl.DomainHandler(w, r, regexp.MustCompile("^/domain/([^/]+)$").FindStringSubmatch(r.URL.Path))
	}))
	defer ts.Close()

	req, err := http.NewRequest("DELETE", ts.URL+"/domain/microsoft.hiv", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer letmein")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify
	d, findErr := repo.FindByName("microsoft.hiv")
	assert.NotNil(findErr)
}
