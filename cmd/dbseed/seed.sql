-- Script to create ledgers, transactions, and tags tables

-- Create the 'ledgers' table
CREATE TABLE ledgers (
    id SERIAL PRIMARY KEY,
    identifier VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    data JSON,
    parent_id INTEGER REFERENCES ledgers(id) ON DELETE SET NULL
);

-- Create the 'transactions' table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_ledger_id INTEGER REFERENCES ledgers(id) ON DELETE CASCADE,
    to_ledger_id INTEGER REFERENCES ledgers(id) ON DELETE CASCADE,
    amount_eur NUMERIC(12, 2) CHECK (amount_eur >= 0), -- Assuming no negative amounts are allowed
    data JSON
);

-- Create the 'tags' table
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INTEGER REFERENCES tags(id) ON DELETE SET NULL
);