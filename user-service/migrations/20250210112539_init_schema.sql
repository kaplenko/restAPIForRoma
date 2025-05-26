-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    pass_hash BYTEA NOT NULL
);

CREATE TABLE orders
(
    id BIGSERIAL PRIMARY KEY,
    order_number TEXT UNIQUE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    accrual INT,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE balances
(
    user_id BIGINT PRIMARY KEY REFERENCES users(id),
    current INT NOT NULL DEFAULT 0,
    withdrawn INT NOT NULL DEFAULT 0
);

CREATE TABLE withdrawals
(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number TEXT NOT NULL,
    sum INT NOT NULL DEFAULT 0,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
