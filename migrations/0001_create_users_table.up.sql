CREATE TABLE users (
    id         VARCHAR(36)  NOT NULL,
    first_name VARCHAR(150) NOT NULL,
    last_name  VARCHAR(150) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at INT          NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
