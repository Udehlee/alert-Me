CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pass_word TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);