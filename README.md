# Click-Counter IFrame

[![Travis](https://travis-ci.org/dothiv/clickcounter-iframe.svg?branch=master)](https://travis-ci.org/dothiv/clickcounter-iframe/)

Simple webservice to host the click-counter iframes via CNAME records.

It is set up as a microservice which is managed by a RESTful API.

## Testing

Create a databse to run the tests on:

    CREATE USER clickcounteriframe;
	CREATE DATABASE clickcounteriframe;
	GRANT ALL PRIVILEGES ON DATABASE clickcounteriframe TO clickcounteriframe;
	
	psql -H localhost -U clickcounteriframe -d clickcounteriframe < sql/domain.sql
