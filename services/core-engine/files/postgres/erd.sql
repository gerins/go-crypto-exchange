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
    quantity                        DOUBLE PRECISION NOT NULL DEFAULT 0,
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

CREATE TABLE pairs (
    id                              SERIAL PRIMARY KEY,
    code                            VARCHAR(256) NOT NULL DEFAULT '',
    primary_crypto_id               INTEGER NOT NULL REFERENCES crypto(id),
    secondary_crypto_id             INTEGER NOT NULL REFERENCES crypto(id),
    status                          BOOLEAN NOT NULL DEFAULT true,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER pairs BEFORE UPDATE ON pairs FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO pairs (
    "id",
    "code",
    "primary_crypto_id",
    "secondary_crypto_id",
    "status",
    "deleted_at"
) VALUES (0, '', 0, 0, false, now());

---------------------------------------------------------------------------------------------------------------------

CREATE TYPE order_side AS ENUM ('BUY', 'SELL');
CREATE TYPE order_type AS ENUM ('MARKET', 'LIMIT', 'STOP_LOSS', 'TAKE_PROFIT')
CREATE TYPE order_status AS ENUM ('COMPLETE', 'FAILED', 'PROGRESS', 'PARTIAL');

CREATE TABLE orders (
    id                              SERIAL PRIMARY KEY,
    user_id                         INTEGER NOT NULL REFERENCES users(id),
    pair_id                         INTEGER NOT NULL REFERENCES pairs(id),
    quantity                        DOUBLE PRECISION NOT NULL DEFAULT 0,
    filled_quantity                 DOUBLE PRECISION NOT NULL DEFAULT 0,
    price                           DOUBLE PRECISION NOT NULL DEFAULT 0,
    type                            order_type,
    side                            order_side,
    status                          order_status,
    transaction_time                BIGINT NOT NULL DEFAULT 0,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER orders BEFORE UPDATE ON orders FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO orders (
    "id",
    "user_id",
    "pair_id",
    "quantity",
    "filled_quantity",
    "price",
    "type",
    "side",
    "status",
    "transaction_time",
    "deleted_at"
) VALUES (0, 0, 0, 0, 0, 0, 'MARKET', 'BUY', 'COMPLETE', 0, now());

---------------------------------------------------------------------------------------------------------------------

CREATE TABLE match_orders (
    id                              SERIAL PRIMARY KEY,
    pair_id                         INTEGER NOT NULL REFERENCES pairs(id),
    taker_order_id                  INTEGER NOT NULL REFERENCES orders(id),
    maker_order_id                  INTEGER NOT NULL REFERENCES orders(id),
    quantity                        DOUBLE PRECISION NOT NULL DEFAULT 0,
    price                           DOUBLE PRECISION NOT NULL DEFAULT 0,
    transaction_time                BIGINT NOT NULL DEFAULT 0,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER match_orders BEFORE UPDATE ON match_orders FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO match_orders (
    "id",
    "pair_id",
    "taker_order_id",
    "maker_order_id",
    "quantity",
    "price",
    "transaction_time",
    "deleted_at"
) VALUES (0, 0, 0, 0, 0, 0, 0, now());
