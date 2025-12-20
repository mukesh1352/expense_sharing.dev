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