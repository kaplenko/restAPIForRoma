-- +goose Up
-- +goose StatementBegin
-- Тестовые данные для таблицы users
INSERT INTO users (username, pass_hash) VALUES
('user1', '$2a$10$LTHwFZUDewwTHBY.1fTRIuJhemIJSHLs69jEFjK1nt5ez/tWltIce'),
('user2', '$2a$10$LTHwFZUDewwTHBY.1fTRIuJhemIJSHLs69jEFjK1nt5ez/tWltIce'), -- password
('user3', '$2a$10$LTHwFZUDewwTHBY.1fTRIuJhemIJSHLs69jEFjK1nt5ez/tWltIce');


-- Тестовые данные для таблицы orders
INSERT INTO orders (order_number, user_id, status, accrual) VALUES
('79927398713', 1, 'PROCESSED', 500),  -- Валидный номер по Луну
('12345678903', 2, 'PROCESSING', NULL), -- Валидный номер
('4242424242424242', 3, 'NEW', NULL);   -- Валидный номер

-- Тестовые данные для таблицы balances
INSERT INTO balances (user_id, current, withdrawn) VALUES
(1, 1000, 200),
(2, 500, 0),
(3, 0, 0);

-- Тестовые данные для таблицы withdrawals
INSERT INTO withdrawals (user_id, order_number, sum) VALUES
(1, '54321', 200),
(2, '67890', 50);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE withdrawals, balances, orders RESTART IDENTITY;
-- +goose StatementEnd
