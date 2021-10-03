CREATE TABLE orders (
    id VARCHAR(36) NOT NULL,
    status VARCHAR(36) NOT NULL,
    payment_method VARCHAR(50),
    payment_intent VARCHAR(50),
    book_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP  NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT orders_pkey PRIMARY KEY (id)
);
