-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, pass_hash) VALUES
('user1', '$2a$10$LTHwFZUDewwTHBY.1fTRIuJhemIJSHLs69jEFjK1nt5ez/tWltIce'),
('user2', '$2a$10$LTHwFZUDewwTHBY.1fTRIuJhemIJSHLs69jEFjK1nt5ez/tWltIce'), -- password
('user3', '$2a$10$LTHwFZUDewwTHBY.1fTRIuJhemIJSHLs69jEFjK1nt5ez/tWltIce');

INSERT INTO orders (order_number, user_id, status, accrual) VALUES
('79927398713', 1, 'PROCESSED', 500),
('12345678903', 2, 'PROCESSING', NULL),
('4242424242424242', 3, 'NEW', NULL);

INSERT INTO balances (user_id, current, withdrawn) VALUES
(1, 10005, 200),
(2, 5005, 0),
(3, 0, 0);

INSERT INTO withdrawals (user_id, order_number, sum) VALUES
(1, '54321', 2002),
(2, '67890', 505);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE withdrawals, balances, orders, users RESTART IDENTITY;
-- +goose StatementEnd
