package clickcounteriframe

import (
	"database/sql"
	"testing"

	"code.google.com/p/gcfg"
	assert "github.com/stretchr/testify/assert"
)

func TestThatItPersistsADomain(t *testing.T) {
	assert := assert.New(t)

	c := NewDefaultConfig()
	configErr := gcfg.ReadFileInto(c, "config.ini")
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())
	db.Exec("TRUNCATE domain RESTART IDENTITY")

	// Persist
	domain := new(Domain)
	domain.Name = "thjnk.hiv"
	domain.Redirect = "http://www.thjnk.de/"
	repo := NewDomainRepository(db)
	persistErr := repo.Persist(domain)
	assert.Nil(persistErr)

	// Verify
	d, findErr := repo.FindByName("thjnk.hiv")
	assert.Nil(findErr)

	assert.Equal(1, d.Id)
	assert.Equal("thjnk.hiv", d.Name)
	assert.Equal("http://www.thjnk.de/", d.Redirect)
}
