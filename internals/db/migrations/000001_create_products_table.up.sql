CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name_ TEXT NOT NULL,
  price DOUBLE PRECISION NOT NULL,
  product_url TEXT NOT NULL,
  status_ TEXT DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT NOW(),
  last_checked TIMESTAMP 
);
