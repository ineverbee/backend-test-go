CREATE TABLE users (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    balance INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id INT,
    balance_change INT NOT NULL,
    comment VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);