-- Table for holders
CREATE TABLE IF NOT EXISTS holders (
    id SERIAL PRIMARY KEY,
    type VARCHAR(31),
    identifier VARCHAR(255),
    name VARCHAR(255),
    parent_holder_id INT,
    data JSONB,
    favorite BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_type_identifier UNIQUE (type, identifier),
    CONSTRAINT fk_parent_holder FOREIGN KEY (parent_holder_id)
        REFERENCES holders (id) ON DELETE SET NULL
);

-- Table for transactions
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_holder_id INT NOT NULL,
    to_holder_id INT NOT NULL,
    amount INT,
    timestamp TIMESTAMP,
    data JSONB,
    parent_transaction_id INT,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_from_holder FOREIGN KEY (from_holder_id)
        REFERENCES holders (id) ON DELETE CASCADE,
    CONSTRAINT fk_to_holder FOREIGN KEY (to_holder_id)
        REFERENCES holders (id) ON DELETE CASCADE,
    CONSTRAINT fk_parent_transaction FOREIGN KEY (parent_transaction_id)
        REFERENCES transactions (id) ON DELETE SET NULL
);

-- Table for tags
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    parent_tag_id INT,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_parent_tag FOREIGN KEY (parent_tag_id)
        REFERENCES tags (id) ON DELETE SET NULL
);

-- Many-to-many relation table between holders and tags
CREATE TABLE IF NOT EXISTS holders_tags (
    holder_id INT NOT NULL,
    tag_id INT NOT NULL,
    CONSTRAINT fk_holder FOREIGN KEY (holder_id)
        REFERENCES holders (id) ON DELETE CASCADE,
    CONSTRAINT fk_tag FOREIGN KEY (tag_id)
        REFERENCES tags (id) ON DELETE CASCADE,
    CONSTRAINT pk_holders_tags PRIMARY KEY (holder_id, tag_id)
);

-- Many-to-many relation table between transactions and tags
CREATE TABLE IF NOT EXISTS transactions_tags (
    transaction_id INT NOT NULL,
    tag_id INT NOT NULL,
    CONSTRAINT fk_transaction FOREIGN KEY (transaction_id)
        REFERENCES transactions (id) ON DELETE CASCADE,
    CONSTRAINT fk_tag_transaction FOREIGN KEY (tag_id)
        REFERENCES tags (id) ON DELETE CASCADE,
    CONSTRAINT pk_transactions_tags PRIMARY KEY (transaction_id, tag_id)
);
