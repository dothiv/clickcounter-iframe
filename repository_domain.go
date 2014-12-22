package clickcounteriframe

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type DomainRepositoryInterface interface {
	Persist(domain *Domain) (err error)
	Remove(domain *Domain) (err error)
	FindByName(name string) (domain *Domain, err error)
}

type DomainRepository struct {
	DomainRepositoryInterface
	db            *sql.DB
	TABLE_NAME    string
	FIELDS        string
	OFFSET_FIELD  string
	CREATED_FIELD string
	UPDATED_FIELD string
	fields        []string
}

func NewDomainRepository(db *sql.DB) (repo *DomainRepository) {
	repo = new(DomainRepository)
	repo.db = db
	repo.TABLE_NAME = "domain"
	repo.FIELDS = "name, redirect, landingpage"
	repo.OFFSET_FIELD = "id"
	repo.CREATED_FIELD = "created"
	repo.UPDATED_FIELD = "updated"
	repo.fields = []string{repo.OFFSET_FIELD, repo.FIELDS, repo.CREATED_FIELD, repo.UPDATED_FIELD}
	return
}

func (repo *DomainRepository) Persist(domain *Domain) (err error) {
	domain.LandingPageJson, err = json.Marshal(domain.LandingPage)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	if domain.Id > 0 {
		_, err = repo.db.Exec("UPDATE "+repo.TABLE_NAME+" "+
			"SET redirect = $1, landingpage = $2 WHERE id = $3", domain.Redirect, domain.LandingPageJson, domain.Id)
	} else {
		err = repo.db.QueryRow("INSERT INTO "+repo.TABLE_NAME+" "+
			"("+repo.FIELDS+") "+
			"VALUES($1, $2, $3) RETURNING id, created",
			domain.Name, domain.Redirect, domain.LandingPageJson).Scan(&domain.Id, &domain.Created)
	}
	return
}

func (repo *DomainRepository) Remove(domain *Domain) (err error) {
	_, err = repo.db.Exec("DELETE FROM "+repo.TABLE_NAME+" "+
		"WHERE "+repo.OFFSET_FIELD+" = $1",
		domain.Id)
	return
}

func (repo *DomainRepository) FindByName(name string) (domain *Domain, err error) {
	domain = new(Domain)

	err = repo.db.QueryRow("SELECT "+strings.Join(repo.fields, ",")+" FROM "+repo.TABLE_NAME+" WHERE name = $1", name).Scan(&domain.Id, &domain.Name, &domain.Redirect, &domain.LandingPageJson, &domain.Created, &domain.Updated)
	err = json.Unmarshal(domain.LandingPageJson, &domain.LandingPage)
	if err != nil {
		return
	}
	return
}
