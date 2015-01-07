DROP TABLE IF EXISTS domain;

CREATE TABLE domain (
	id SERIAL PRIMARY KEY NOT NULL UNIQUE,
	name varchar(128) NOT NULL UNIQUE,
	redirect TEXT NULL DEFAULT NULL,
	landingpage JSON NULL DEFAULT NULL,
	created timestamp DEFAULT current_timestamp,
	updated timestamp DEFAULT current_timestamp
);

CREATE INDEX domain__dn_idx ON domain ( name );
