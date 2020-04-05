CREATE TABLE IF NOT EXISTS client (
	id           TEXT  NOT NULL PRIMARY KEY,
	secret 		 TEXT  NOT NULL,
	domain       TEXT  NOT NULL
);
