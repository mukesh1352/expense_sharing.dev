
CREATE TABLE settlements (
    id UUID PRIMARY KEY,
    from_user_id UUID REFERENCES users(id) NOT NULL,
    to_user_id UUID REFERENCES users(id) NOT NULL,
    amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT NOW(),
    CHECK (from_user_id <> to_user_id)
);
