![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)

# Expense Sharing System

This project is a **backend-first service** that implements a centralized
expense-sharing ledger inspired by Splitwise.

It allows users to create groups, add shared expenses, track who owes whom,
and settle outstanding dues. The primary focus of this project is **financial
correctness, balance simplification, and system design**, rather than UI
complexity.

---

## Problem Statement

The goal of this system is to design a backend service that enables:

- Creating groups for users
- Adding shared expenses within a group
- Supporting multiple split types (Equal, Exact, Percentage)
- Tracking balances to determine who owes whom
- Simplifying balances to minimize the number of transactions
- Settling outstanding dues

The system applies **distributed system design principles** while maintaining
a **centralized, strongly consistent ledger** for financial correctness.

---

## Tech Stack

- **Backend**: Golang
- **Database**: PostgreSQL (NEON)
- **Frontend**: React
- **API Style**: REST

---

## Design Choices

- **PostgreSQL** is used to provide strong consistency and transactional
  guarantees, which are essential for financial data. **Neon** is preffered in this project.
- The backend is designed as a **stateless service**, with the database acting
  as the **single source of truth**.
- Authentication and payment processing are intentionally kept **out of scope**
  to focus on ledger correctness and balance management.

---

## Data Model Overview

The core entities in the system are:

- User
- Group
- Expense
- Expense Split
- Balance
- Settlement

Balances are stored in a **directional format**:

```

from_user_id → to_user_id → amount

```

This represents **who owes whom and how much**.

The system follows a **ledger-based, centralized model**.

---

## What Is a Ledger?

A ledger is a system that:

- Records who owes whom
- Stores net financial obligations
- Preserves a history of settlements
- Guarantees correctness and consistency
- Acts as a single source of truth

This project implements a **financial obligation ledger**, not a payment system.

---

## Ledger Architecture

The system is designed as a centralized financial ledger where different tables
serve distinct roles in tracking obligations and settlements.

| Table            | Ledger Role                               |
|------------------|--------------------------------------------|
| `expenses`       | Events that create financial obligations   |
| `balances`       | Current state of net obligations            |
| `settlements`    | Immutable historical records of settlements |
| `expense_splits` | Defines how obligations are derived         |
| `group_members`  | Validates user participation in a group     |

---

## Balance Model

The system maintains a **directional balance ledger** where each entry
represents a net obligation between two users.

A balance entry of the form:

```

from_user_id → to_user_id = amount

```

means that `from_user_id` owes `to_user_id` the specified amount.

### Balance Invariants

The following rules are always enforced:

- Balances are always positive
- A user never owes themselves
- Only one directional balance exists between any two users
- Balances represent net obligations across all expenses

---

## Mandatory Balance Constraints

- `amount > 0`
- `from_user_id ≠ to_user_id`
- Unidirectional flow between any two users
- Duplicate balances are not allowed

---

## Expense Processing Flow

When an expense is added, the system **does not move money**. Instead, it
creates financial obligations between users.

The general flow is:

1. One user pays the full expense amount (the payer).
2. Each participant owes a calculated share based on the selected split type.
3. For each participant (excluding the payer), a balance is updated to reflect
   that the participant owes the payer their share.
4. Balances are netted and simplified to remove redundant obligations.
5. All operations occur within a **single database transaction**.

---

## Supported Split Types

### Equal Split
The total expense amount is divided equally among all participants.

### Exact Amount Split
Each participant owes a specific amount.

```

sum(participant_shares) = total_amount

```

### Percentage Split
Each participant owes a percentage of the total expense.

```

sum(percentages) = 100%

```

---

## Balance Netting and Simplification

The system stores balances as **net financial obligations** to keep the ledger
minimal and easy to understand.

### Balance Netting

If two users owe each other, the system nets the balances so that only the
difference is stored.

**Example:**

- User A owes User B ₹100
- User B owes User A ₹60

After netting:

- User A owes User B ₹40

This avoids redundant or contradictory debts.

---

### Balance Simplification

Beyond pairwise netting, the system also simplifies **multi-party debts**.

If:
- User A owes User B
- User B owes User C

Then the system simplifies this into a **direct obligation from User A to
User C**, wherever possible.

**Example:**

- User A → User B = ₹100  
- User B → User C = ₹100  

Simplified to:

- User A → User C = ₹100  

This minimizes the number of outstanding balances and reduces the number of
payments required to settle all dues, while preserving the correct net outcome.

---

## Settlements

Settlements represent the fulfillment of an obligation **outside the system**,
such as through cash, bank transfer, or UPI.

When a settlement is recorded:

- The corresponding balance is reduced or removed
- Fully settled balances are deleted from the ledger

All settlements are stored as **immutable records** to provide an auditable
history of payments, while balances always reflect the **current outstanding
obligations**.

---
