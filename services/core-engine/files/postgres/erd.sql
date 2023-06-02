SET TIMEZONE = 'Etc/GMT-7';

CREATE OR REPLACE FUNCTION update_modified_column() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW; 
END;
$$ language 'plpgsql';

---------------------------------------------------------------------------------------------------------------------

CREATE TABLE users (
    id                              SERIAL PRIMARY KEY,
    full_name                       VARCHAR(128) NOT NULL DEFAULT '',
    email                           VARCHAR(128) NOT NULL,
    phone_number                    VARCHAR(128) NOT NULL DEFAULT '',
    balance                         BIGINT NOT NULL DEFAULT 0,
    password                        VARCHAR(512) NOT NULL,
    status                          BOOLEAN NOT NULL DEFAULT true,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER users BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO users (
    "id",
    "full_name",
    "email",
    "phone_number",
    "password",
    "status",
    "deleted_at"
) VALUES (0, '', '', '', '', false, now());

---------------------------------------------------------------------------------------------------------------------

CREATE TABLE crypto (
    id                              SERIAL PRIMARY KEY,
    symbol                          VARCHAR(128) NOT NULL DEFAULT '', -- BTC, ETH
    status                          BOOLEAN NOT NULL DEFAULT true,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER crypto BEFORE UPDATE ON crypto FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO crypto (
    "id",
    "symbol",
    "status",
    "deleted_at"
) VALUES (0, '', false, now());

---------------------------------------------------------------------------------------------------------------------

CREATE TABLE wallet (
    id                              SERIAL PRIMARY KEY,
    user_id                         INTEGER NOT NULL REFERENCES users(id),
    crypto_id                       INTEGER NOT NULL REFERENCES crypto(id),
    quantity                        INTEGER NOT NULL DEFAULT 0,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER wallet BEFORE UPDATE ON wallet FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO wallet (
    "id",
    "user_id",
    "crypto_id",
    "quantity",
    "deleted_at"
) VALUES (0, 0, 0, 0, now());

---------------------------------------------------------------------------------------------------------------------

CREATE TYPE order_type AS ENUM ('BUY', 'SELL');
CREATE TYPE order_status AS ENUM ('COMPLETE', 'FAILED', 'PROGRESS', 'PARTIAL');

CREATE TABLE orders (
    id                              SERIAL PRIMARY KEY,
    crypto_id                       INTEGER NOT NULL REFERENCES crypto(id),
    user_id                         INTEGER NOT NULL REFERENCES users(id),
    quantity                        INTEGER NOT NULL DEFAULT 0,
    price                           INTEGER NOT NULL DEFAULT 0,
    type                            order_type,
    status                          order_status,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER orders BEFORE UPDATE ON orders FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO orders (
    "id",
    "crypto_id",
    "user_id",
    "quantity",
    "price",
    "type",
    "status",
    "deleted_at"
) VALUES (0, 0, 0, 0, 0, 'BUY', 'COMPLETE', now());

---------------------------------------------------------------------------------------------------------------------

CREATE TABLE match_orders (
    id                              SERIAL PRIMARY KEY,
    crypto_id                       INTEGER NOT NULL REFERENCES crypto(id),
    buy_order_id                    INTEGER NOT NULL REFERENCES orders(id),
    sell_order_id                   INTEGER NOT NULL REFERENCES orders(id),
    quantity                        INTEGER NOT NULL DEFAULT 0,
    price                           INTEGER NOT NULL DEFAULT 0,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER match_orders BEFORE UPDATE ON match_orders FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO match_orders (
    "id",
    "crypto_id",
    "buy_order_id",
    "sell_order_id",
    "quantity",
    "price",
    "deleted_at"
) VALUES (0, 0, 0, 0, 0, 0, now());
