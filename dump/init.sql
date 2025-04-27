CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    reference VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    price int NOT NULL,
    image_url TEXT NOT NULL,
    product_url TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    event_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    in_stock BOOLEAN DEFAULT true
);

CREATE TABLE scrape_targets (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    category VARCHAR(100) NOT NULL,
    active BOOLEAN DEFAULT true
);


CREATE INDEX idx_products_product_id ON products(reference);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_event_type ON products(event_type);
CREATE INDEX idx_products_date ON products(event_date);

CREATE INDEX idx_scrape_targets_url ON scrape_targets(url);