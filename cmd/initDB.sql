-- create database manualy before this script

CREATE TABLE IF NOT EXISTS users (
    id INT NOT NULL PRIMARY KEY,
    login TEXT NOT NULL DEFAULT '',
    password TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS session (
    user_id INT NOT NULL,
    session TEXT NOT NULL DEFAULT '' UNIQUE
)