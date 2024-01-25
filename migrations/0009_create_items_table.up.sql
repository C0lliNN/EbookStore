CREATE TABLE items
(
    id               VARCHAR(36)  NOT NULL,
    order_id         VARCHAR(36)  NOT NULL,
    name             VARCHAR(255) NOT NULL,
    price            INT          NOT NULL,
    preview_image_id VARCHAR(36)  NOT NULL,
    CONSTRAINT items_pkey PRIMARY KEY (id, order_id),
    CONSTRAINT items_order_id_fkey FOREIGN KEY (order_id)
        REFERENCES orders (id) ON DELETE CASCADE ON UPDATE CASCADE
);
