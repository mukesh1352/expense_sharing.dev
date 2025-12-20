CREATE TABLE expenses (
    id UUID PRIMARY KEY,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    paid_by UUID REFERENCES users(id) NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    split_type TEXT NOT NULL CHECK (split_type IN ('EQUAL', 'EXACT', 'PERCENT')),
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);



