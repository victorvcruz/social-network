CREATE TABLE account (
	id VARCHAR,
	username VARCHAR UNIQUE NOT NULL,
	name VARCHAR UNIQUE NOT NULL,
	description VARCHAR NOT NULL,
	email VARCHAR UNIQUE NOT NULL,
	password VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted BOOLEAN NOT NULL,
        PRIMARY KEY (id)
);

CREATE TABLE account_follow (
	account_id VARCHAR NOT NULL,
	account_id_followed VARCHAR NOT NULL,
	FOREIGN KEY (account_id) REFERENCES account (id),
	FOREIGN KEY (account_id_followed) REFERENCES account (id)
);

CREATE TABLE post (
	id VARCHAR,
	account_id VARCHAR NOT NULL,
	content VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        removed BOOLEAN NOT NULL,
        PRIMARY KEY (id),
        FOREIGN KEY (account_id) REFERENCES account(id)
);

CREATE TABLE comment (
	id VARCHAR,
	account_id VARCHAR NOT NULL,
	post_id VARCHAR NOT NULL,
	content VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        removed BOOLEAN NOT NULL,
        PRIMARY KEY (id),
        FOREIGN KEY (account_id) REFERENCES account (id),
        FOREIGN KEY (post_id) REFERENCES post (id)
);

CREATE TABLE interactions (
	id VARCHAR,
	account_id VARCHAR NOT NULL,
	post_id VARCHAR NOT NULL,
	comment_id VARCHAR,
	type VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        removed BOOLEAN NOT NULL,
        PRIMARY KEY (id),
        FOREIGN KEY (account_id) REFERENCES account (id),
        FOREIGN KEY (post_id) REFERENCES post (id),
        FOREIGN KEY (comment_id) REFERENCES comment (id)
);
