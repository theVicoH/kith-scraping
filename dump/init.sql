CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    reference VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    price int NOT NULL,
    image_url TEXT NOT NULL,
    product_url TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    event_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_product_id ON products(reference);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_event_type ON products(event_type);
CREATE INDEX idx_products_date ON products(event_date);