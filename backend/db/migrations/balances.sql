CREATE TABLE balances (
    from_user_id UUID REFERENCES users(id),
    to_user_id UUID REFERENCES users(id),
    amount NUMERIC(12, 2) NOT NULL CHECK (amount >= 0),
    CHECK (from_user_id <> to_user_id),
    PRIMARY KEY (from_user_id, to_user_id)
);

