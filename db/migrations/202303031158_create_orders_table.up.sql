CREATE TABLE orders
(
    id          STRING(36) NOT NULL,
    customer_id STRING(36) NOT NULL,
    cart_id     STRING(36) NOT NULL,
    fare        NUMERIC   NOT NULL,
    currency    STRING(3) NOT NULL,
    status      STRING(64) NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
) PRIMARY KEY (id),
ROW DELETION POLICY (OLDER_THAN(created_at, INTERVAL 15 DAY));

CREATE UNIQUE INDEX idx_cart_id ON orders (cart_id);
