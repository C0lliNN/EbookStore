CREATE TABLE orders (
    id VARCHAR(36) NOT NULL,
    status VARCHAR(36) NOT NULL,
    total INT NOT NULL,
    payment_intent_id VARCHAR(50),
    client_secret VARCHAR(255),
    book_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP  NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT orders_pkey PRIMARY KEY (id)
);
