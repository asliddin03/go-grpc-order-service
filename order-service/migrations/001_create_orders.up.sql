CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status TEXT NOT NULL,
    total_price BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_up TIMESTAMPTZ NOT NULL
);