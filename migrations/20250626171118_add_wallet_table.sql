-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    walletid BIGINT NOT NULL UNIQUE,
    balance INTEGER NOT NULL DEFAULT 0 CHECK (balance >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wallets;
-- +goose StatementEnd
