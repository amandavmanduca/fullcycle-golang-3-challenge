-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
