CREATE TABLE IF NOT EXISTS delegation (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    tx_hash VARCHAR(255) NOT NULL,
    amount DECIMAL NOT NULL,
    delegator VARCHAR(255) NOT NULL,
    height BIGINT NOT NULL,
    datetime TIMESTAMPTZ NOT NULL,
    CONSTRAINT uq_tx_hash UNIQUE (tx_hash)
);

CREATE TABLE IF NOT EXISTS polling (
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    last_polled_at TIMESTAMPTZ NOT NULL
);