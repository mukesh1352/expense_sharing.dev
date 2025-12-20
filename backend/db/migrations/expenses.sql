CREATE TABLE expenses (
    id UUID PRIMARY KEY,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    paid_by UUID REFERENCES users(id) NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    split_type TEXT NOT NULL CHECK (split_type IN ('EQUAL', 'EXACT', 'PERCENT')),
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE expense_splits (
    expense_id UUID REFERENCES expenses(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    amount NUMERIC(12, 2),
    percentage NUMERIC(5, 2),
    CHECK (
        amount IS NOT NULL OR percentage IS NOT NULL
    ),
    PRIMARY KEY (expense_id, user_id)
);

