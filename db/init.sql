CREATE TABLE account (
	id VARCHAR UNIQUE NOT NULL,
	username VARCHAR NOT NULL,
	name VARCHAR NOT NULL,
	description VARCHAR NOT NULL,
	email VARCHAR NOT NULL,
	password VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted BOOLEAN NOT NULL,
        PRIMARY KEY (id)
);
CREATE UNIQUE INDEX account_username_index ON account(username, deleted) WHERE deleted = false;
CREATE UNIQUE INDEX account_email_index ON account(email, deleted) WHERE deleted = false;

CREATE TABLE account_follow (
	account_id VARCHAR NOT NULL,
	account_id_followed VARCHAR NOT NULL,
	FOREIGN KEY (account_id) REFERENCES account (id),
	FOREIGN KEY (account_id_followed) REFERENCES account (id)
);

CREATE TABLE post (
	id VARCHAR UNIQUE NOT NULL,
	account_id VARCHAR NOT NULL,
	content VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        removed BOOLEAN NOT NULL,
        PRIMARY KEY (id),
        FOREIGN KEY (account_id) REFERENCES account(id)
);

CREATE TABLE comment (
	id VARCHAR UNIQUE NOT NULL,
	account_id VARCHAR NOT NULL,
	post_id VARCHAR NOT NULL,
	comment_id VARCHAR,
	content VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        removed BOOLEAN NOT NULL,
        PRIMARY KEY (id),
        FOREIGN KEY (account_id) REFERENCES account (id),
        FOREIGN KEY (post_id) REFERENCES post (id)
);

CREATE TABLE interaction (
	id VARCHAR UNIQUE NOT NULL,
	account_id VARCHAR NOT NULL,
	post_id VARCHAR NULL,
	comment_id VARCHAR NULL,
	type VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        removed BOOLEAN NOT NULL,
        PRIMARY KEY (id),
        FOREIGN KEY (account_id) REFERENCES account (id),
        FOREIGN KEY (post_id) REFERENCES post (id),
        FOREIGN KEY (comment_id) REFERENCES comment (id)
);
