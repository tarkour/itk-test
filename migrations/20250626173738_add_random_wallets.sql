-- +goose Up
-- +goose StatementBegin
INSERT INTO wallets (walletid, balance) 
VALUES 
  (1001, 500),
  (1002, 1500),
  (1003, 0),
  (1004, 750),
  (1005, 3200),
  (1006, 42),
  (1007, 100000),
  (1008, 999),
  (1009, 1),
  (1010, 128);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM wallets;
-- +goose StatementEnd
