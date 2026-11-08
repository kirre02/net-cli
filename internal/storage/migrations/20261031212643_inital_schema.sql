-- +goose Up

CREATE TABLE IF NOT EXISTS collections (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS requests (
    id TEXT PRIMARY KEY,
    collection_id TEXT NOT NULL,
    name TEXT NOT NULL,
    method TEXT NOT NULL,
    url TEXT NOT NULL,
    headers TEXT NOT NULL DEFAULT '{}',
    body TEXT NOT NULL DEFAULT '',

    FOREIGN KEY (collection_id)
        REFERENCES collections(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_requests_collection_id
    ON requests(collection_id);


-- +goose Down

DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS collections;
