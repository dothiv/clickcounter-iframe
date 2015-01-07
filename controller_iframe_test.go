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

	repo := NewDomainRepository(db)

	thjnk := new(Domain)
	thjnk.Name = "thjnk.hiv"
	thjnk.Redirect.String = "http://www.thjnk.de/"
	thjnk.Redirect.Valid = true
	assert.Nil(repo.Persist(thjnk))

	caro4life := new(Domain)
	caro4life.Name = "caro4life.hiv"
	lp := new(LandingPage)
	lp.DefaultLocale = "de"
	lp.Strings = make(map[string]*LandingPageText)
	lp.Strings["de"] = new(LandingPageText)
	lp.Strings["de"].Locale = "de"
	lp.Strings["de"].Title = "Carolin's digital Red Ribbon"
	lp.Strings["de"].About = "I support the global, digital movement to see the end of AIDS. Together we can do it! Let’s build an AIDS free generation."
	lp.Strings["de"].MicroDonation = "Every visit to this website triggers a much needed donation to global HIV projects. Thank you!"
	lp.Strings["de"].LearnMore = "Learn more"
	lp.Strings["de"].GetYourOwn = "Get your own"
	lp.Strings["de"].TellYourFriends = "Tell your friends:"
	lp.Strings["de"].Tweet = "Lets work together to see an #AIDS free generation @dotHIV – caro4life.hiv"
	lp.Strings["de"].Imprint = "Imprint"

	caro4life.LandingPage = lp
	assert.Nil(repo.Persist(caro4life))

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

func TestThatItReturnsTheLandingpage(t *testing.T) {
	assert := assert.New(t)

	cntrl := SetupDomainTest(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = "caro4life.hiv"
		r.Header.Add("Accept-Language", "de-DE,de;q=0.8,en;q=0.6")
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
	assert.Contains(body, `<html lang="de"`)
	assert.Contains(body, `Carolin's digital Red Ribbon`)
}
