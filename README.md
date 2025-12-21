![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![TypeScript](https://img.shields.io/badge/typescript-%23007ACC.svg?style=for-the-badge&logo=typescript&logoColor=white)
![TailwindCSS](https://img.shields.io/badge/tailwindcss-%2338B2AC.svg?style=for-the-badge&logo=tailwind-css&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgresql-%23336791.svg?style=for-the-badge&logo=postgresql&logoColor=white)


# Expense Sharing System

This project is a **backend-first service** that implements a centralized
expense-sharing ledger.

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
- **Database**: PostgreSQL (Supabase)
- **Frontend**: React
- **API Style**: REST

---

## Design Choices

- **PostgreSQL** is used to provide strong consistency and transactional
  guarantees, which are essential for financial data. **Supabase** is used as the managed PostgreSQL provider for reliability and always-on availability during development.
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

## Database Schema Management
The database schema is managed using the SQL migration files which are located in the `backend/db/migrations` directory.

Each of the migration files defines a particular part of the database schema, such as the tables, its respective relations and the constraints required by the ledger system.

Migrations are applied **outside of the application runtime** and are not
executed by the backend service itself.
This enures that :
 - The application does not modify the schema in production.
 - Schema changes are modified and auditable.
 - No race conditions occur.
 - For clear understanding of the data model and the architecture.

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
## Core Backend Functions and Responsibilities

The backend follows a **ledger-centric design**, where all financial state
changes are performed through a small set of well-defined, transactional
functions. Each function has a single responsibility and preserves ledger
correctness at all times.

---

### Transaction Management

#### `withTx(fn)`

Executes a sequence of database operations within a single transaction.

**Responsibilities:**
- Begins a database transaction
- Commits on success
- Rolls back automatically on error
- Helps in the denial of the race condition by using the rollback and commit protocols.

**Why it exists:**
All financial operations must be **atomic**. Partial updates could corrupt
balances and violate ledger invariants.

---

### Expense Creation

#### `CreateExpense(ctx, input)`

Creates a new expense and updates the ledger accordingly.

**Executed inside one transaction.**

**Responsibilities:**
1. Validates expense input
2. Inserts the expense record
3. Calculates how much each participant owes
4. Inserts `expense_splits`
5. Updates balances using ledger core logic
6. Commits atomically


---

### Split Calculation

#### `calculateShares(input)`

Determines how the total expense amount is divided among participants.

Delegates to one of the following based on `split_type`:

- `calculateEqualSplit`
- `calculateExactSplit`
- `calculatePercentageSplit`

These functions are **pure business logic** and do not interact with the database.

---

### Balance Updates

#### `applyBalanceDelta(tx, fromUser, toUser, amount)`

Applies a single financial obligation to the ledger.

**Responsibilities:**
- Prevents self-debt
- Enforces positive balances
- Detects reverse balances
- Performs netting when required
- Ensures only one directional balance exists between users

This function is the **core of ledger correctness**.

---

### Balance Simplification

#### `SimplifyUserBalances(tx, userID)`

Simplifies balances involving a specific user by collapsing indirect obligations.

**Example:**
- A owes B
- B owes C  
→ simplified to  
- A owes C

---

#### `SimplifyBalances(tx)`

Runs balance simplification across all users to keep the ledger minimal and
efficient after expense creation.

---

### Settlement Processing

#### `SettleBalance(ctx, fromUser, toUser, amount)`

Records a real-world payment and updates the ledger.

**Responsibilities:**
- Validates settlement amount
- Inserts an immutable settlement record
- Reduces or removes the corresponding balance

Settlements represent payments **outside the system** (cash, bank transfer, UPI).

---

### Read Operations

#### `GetBalancesForUser(userID)`

Returns all balances where the user is either:
- Owing money, or
- Being owed money

Used to display:
- “You owe”
- “You are owed”

---

#### `GetGroupBalances(groupID)`

Returns all balances between members of a specific group.


## Balance Direction, Netting, and Ledger Terminology

The system maintains balances as **directional financial obligations** between users.
Each balance represents a net amount that one user owes another.

A balance entry is stored as:

from_user_id → to_user_id = amount

This means that `from_user_id` owes `to_user_id` the specified amount.

---

### Forward Balance

A **forward balance** exists when a new obligation is applied in the same direction
as an existing balance.

**Example:**

Existing balance:
```

Bob → Alice = 30

```

New obligation:
```

Bob → Alice = 50

```

Resulting balance:
```

Bob → Alice = 80

```

In this case, the amounts are simply added because the direction is the same.

---

### Reverse Balance

A **reverse balance** exists when a new obligation is applied in the opposite direction
of an existing balance.

**Example:**

Existing balance:
```

Alice → Bob = 30

```

New obligation:
```

Bob → Alice = 50

```

These balances represent opposite obligations and must be **netted** to avoid
redundant or contradictory debts.

---

### Balance Netting

**Netting** is the process of canceling out obligations in opposite directions and
keeping only the net amount.

#### Case 1: Reverse balance is greater
```

Alice → Bob = 100
Bob → Alice = 40

```

Net result:
```

Alice → Bob = 60

```

#### Case 2: Reverse balance is smaller
```

Alice → Bob = 30
Bob → Alice = 80

```

Net result:
```

Bob → Alice = 50

```

#### Case 3: Reverse balance is equal
```

Alice → Bob = 50
Bob → Alice = 50

```

Net result:
```

(no balance exists)

```

Zero balances are removed from the ledger, as they represent no outstanding
obligation.

---

### Why Reverse Balances Are Always Checked First

When applying a new obligation, the system always checks for a **reverse balance**
before updating a forward balance.

This ensures:
- Only one balance exists between any two users
- The ledger stores net obligations only
- Circular or duplicate debts are prevented

Allowing both directions to exist would make it impossible to clearly determine
how much a user actually owes.

---

### Ledger Invariants Enforced

At all times, the ledger enforces the following rules:

- Balances are always positive
- A user never owes themselves
- Only one directional balance exists between any two users
- Reverse balances are netted before forward balances are created or updated
- Zero-value balances are removed

These invariants guarantee that the ledger remains minimal, consistent, and
financially correct.

---

### Key Takeaway

Every time a new obligation is added, the system:

1. Cancels any existing obligation in the opposite direction
2. Applies only the remaining net amount
3. Ensures the ledger reflects the true financial state

## API Overview

| Method | Endpoint            | Description                          |
|------|---------------------|--------------------------------------|
| GET  | `/balances/user`    | Get balances for a user              |
| GET  | `/balances/groups`  | Get balances within a group          |
| POST | `/expenses`         | Create a new expense                 |
| POST | `/settle`           | Record a settlement                  |

---

## Running the Project Locally

This project is designed as a **backend-first system** with a simple frontend
for demonstration purposes. The backend and frontend are run as separate
services.

---

### Prerequisites

Ensure the following tools are installed:

- Go **1.22+**
- Node.js **18+**
- PostgreSQL (or Supabase account)
- Docker (optional, for containerized backend)

---

### Clone the Repository

```bash
git clone https://github.com/<your-username>/<your-repo-name>.git
cd <your-repo-name>
````

---

## Backend Setup

### 1. Configure Environment Variables

Create a `.env` file inside the `backend` directory:

```env
DATABASE_URL=postgresql://<username>:<password>@<host>:<port>/<database>?sslmode=require
```

> **Note:**
> The backend expects an existing PostgreSQL database.
> In development, Supabase is used as a managed PostgreSQL provider.

---

### 2. Apply Database Migrations

Database schema is managed using SQL migration files located at:

```
backend/db/migrations
```

Migrations should be applied **manually** using a PostgreSQL client:

```bash
psql "$DATABASE_URL" -f backend/db/migrations/users.sql
psql "$DATABASE_URL" -f backend/db/migrations/groups.sql
psql "$DATABASE_URL" -f backend/db/migrations/group_members.sql
psql "$DATABASE_URL" -f backend/db/migrations/expenses.sql
psql "$DATABASE_URL" -f backend/db/migrations/expense_splits.sql
psql "$DATABASE_URL" -f backend/db/migrations/balances.sql
psql "$DATABASE_URL" -f backend/db/migrations/settlements.sql
```

---

### 3. (Optional) Seed Sample Data

To populate the database with sample users and groups:

```bash
psql "$DATABASE_URL" -f backend/db/seed.sql
```

---

### 4. Run Backend Locally

```bash
cd backend
go run main.go
```

The backend server will start on:

```
http://localhost:8080
```

---

## Frontend Setup

### 1. Install Dependencies

```bash
cd expense-sharing
pnpm install
```

---

### 2. Configure Frontend Environment Variables

Create a `.env` file inside the `expense-sharing` directory:

```env
VITE_API_BASE=http://localhost:8080
```

---

### 3. Run Frontend Locally

```bash
pnpm run dev
```

The frontend will be available at:

```
http://localhost:5173
```

---

## Running with Docker (Optional)

The backend can also be run using Docker.

### Build and Run Backend Container

```bash
cd backend
docker build -t expense-ledger .
docker run -p 8080:8080 --env-file .env expense-ledger
```

---

## Notes

* Authentication and authorization are intentionally out of scope.
* All financial state changes occur through database transactions.
* The database acts as the single source of truth.
* Frontend auto-refresh mechanisms are intentionally minimal to keep focus on
  ledger correctness and system design.

---

## Key Focus of the Project

This project prioritizes:

* Ledger correctness over UI complexity
* Strong consistency and transactional safety
* Clear domain modeling
* Simplified, net financial obligations
