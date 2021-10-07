CREATE TABLE books (
    id VARCHAR(36) NOT NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    author_name VARCHAR(100) NOT NULL,
    poster_image_bucket_key VARCHAR(50) NOT NULL,
    content_bucket_key VARCHAR(50) NOT NULL,
    release_date DATE NULL,
    price INTEGER NOT NULL,
    created_at TIMESTAMP  NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT books_pkey PRIMARY KEY (id)
);
