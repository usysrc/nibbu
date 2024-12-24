-- init.sql

-- items table
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR NOT NULL
);

-- users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- seed the db
INSERT INTO items (name) SELECT 'item1' WHERE NOT EXISTS (SELECT 1 FROM items WHERE name = 'item1');
-- INSERT INTO users (username,password) VALUES ('test', 'p@ssword')