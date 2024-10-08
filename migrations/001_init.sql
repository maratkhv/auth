-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    login TEXT UNIQUE,
    password BYTEA NOT NULL
);