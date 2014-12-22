package clickcounteriframe

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupDomainTest(t *testing.T) (cntrl *IframeController) {
	assert := assert.New(t)
	c, configErr := NewConfig("config.ini")
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())

	domainRepo := NewDomainRepository(db)
	cntrl = NewIframeController(domainRepo)
	db.Exec("TRUNCATE domain RESTART IDENTITY")

	domain := new(Domain)
	domain.Name = "thjnk.hiv"
	domain.Redirect = "http://www.thjnk.de/"
	repo := NewDomainRepository(db)
	persistErr := repo.Persist(domain)
	assert.Nil(persistErr)

	return
}

func TestThatItReturnsTheIframe(t *testing.T) {
	assert := assert.New(t)

	cntrl := SetupDomainTest(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = "thjnk.hiv"
		cntrl.IframeHandler(w, r, []string{})
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	body := string(b)
	assert.Contains(res.Header.Get("Content-Type"), "text/html")
	assert.Contains(body, "<title>thjnk.hiv</title>")
	assert.Contains(body, `<iframe src="http://www.thjnk.de/" width="100%" height="100%" id="clickcounter-target-iframe"></iframe>`)
}
