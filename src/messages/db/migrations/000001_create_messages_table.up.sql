CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender VARCHAR(70) NOT NULL,
    recipient VARCHAR(70) NOT NULL,
    content TEXT CHECK (LENGTH(content) <= 1024) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    shard_key TEXT NOT NULL
);
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_pkey;
SELECT create_distributed_table('messages', 'shard_key');