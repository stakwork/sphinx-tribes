CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    event VARCHAR(50) NOT NULL,
    pubkey VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    retries INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'PENDING',
    uuid UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
