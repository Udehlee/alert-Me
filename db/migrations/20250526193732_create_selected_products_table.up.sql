  CREATE TABLE selectedProducts (
    id UUID SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    product_id TEXT NOT NULL,
    product_name TEXT NOT NULL,
    current_price DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
    );