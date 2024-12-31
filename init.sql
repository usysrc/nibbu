-- init.sql

-- user table
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- post table
CREATE TABLE IF NOT EXISTS post (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    VARCHAR NOT NULL,
    content VARCHAR NOT NULL,
    url     VARCHAR NOT NULL,
    author  INTEGER,
    date DATETIME DEFAULT CURRENT_TIMESTAMP,
    published VARCHAR,
    FOREIGN KEY(author) REFERENCES user(id)
);

-- seed the db
-- INSERT into user (username,password) VALUES ('test', 'vs');
-- INSERT into post (name,content,url,author,date) VALUES ('my first blog post','hello world', 'my-first-blog-post', 1, "2021-12-09T16:34:04Z");