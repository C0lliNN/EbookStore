CREATE TABLE images
(
    id                      VARCHAR(100)  NOT NULL,
    description             TEXT         NOT NULL,
    book_id            VARCHAR(36),
    CONSTRAINT images_pkey PRIMARY KEY (id),
    CONSTRAINT book_id_fk FOREIGN KEY (book_id) REFERENCES books (id) 
    ON DELETE CASCADE ON UPDATE CASCADE
);
