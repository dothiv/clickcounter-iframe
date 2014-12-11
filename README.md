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

Run the tests

    go test ./...

## About

This web-microservice holds (in a PostgreSQL database) the configuration for our .hiv domains. At the moments that is only the redirect url. .hiv domains can be configured to use `iframe.clickcounter.hiv` as the CNAME record for their `www` subdomain. The webservice extracts the hiv domain name from the host and uses this to render the iframe ([see test](https://github.com/dothiv/clickcounter-iframe/blob/master/controller_iframe_test.go)).

The hiv domain configuration is [updated by our portal](https://github.com/dothiv/dothiv/commit/a8185b8e905b8f4e060c16de919b7f7ef68958a4) immidiately after a user changes to configuration for his .hiv domain by calling the `/domain/$domain` endpoints of the api. This endpoints supports creating, updating and deleting the domain configuration ([see test](https://github.com/dothiv/clickcounter-iframe/blob/master/controller_admin_test.go)).

This is an exemplary curl call for adding / updating a domain configuration:

    curl -u iframe-admin:changeme -X PUT -H "Content-type: application/json" \
    -d '{"redirect":"http://thjnk.de"}' http://iframe.clickcounter.hiv/domain/thjnk.hiv

The webservice is hosted on a Google Cloud Compute Micro instance. It is run behind an nginx who manages the authentication of the admin endpoints.

	# cat /etc/nginx/sites-enabled/iframe 
	server {
	        listen          80;
	        server_name     192.158.31.86 *.hiv;
	        location /domain {
	                auth_basic "Click-counter Iframe Admin Area";
	                auth_basic_user_file /etc/nginx/.htpasswd;
	                proxy_pass http://127.0.0.1:8887/domain;
	                proxy_set_header Host            $host;
	                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	        }
	        location / {
	                proxy_pass http://127.0.0.1:8887;
	                proxy_set_header Host            $host;
	                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	        }
	}

The webservice process is monitored by *upstart*.

	# cat /etc/init/clickcounter-iframe.conf 
	description     "click-counter iframe"
	author          "Markus Tacker <m@tld.hiv>"
	start on (net-device-up
	          and local-filesystems
	          and runlevel [2345])
	stop on runlevel [06]
	respawn
	script
	    set -x
	    chdir /var/www/iframe
	    exec sudo -u iframe /var/www/iframe/clickcounter-iframe /var/www/iframe/config.ini
	end script
