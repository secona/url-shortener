CREATE TABLE links (
	slug TEXT PRIMARY KEY NOT NULL,
	link TEXT NOT NULL,
	user_id INTEGER,

	FOREIGN KEY (user_id) REFERENCES users(id)
);
