-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS states (
    key varchar(255) NOT NULL PRIMARY KEY,
    value varchar(255) NOT NULL UNIQUE
);
INSERT INTO states (key, value)
VALUES ('IsMaintenance', false);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS states;
-- +goose StatementEnd
