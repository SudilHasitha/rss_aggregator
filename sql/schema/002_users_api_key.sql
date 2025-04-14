CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE users 
ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
    encode(
        digest(
            random()::text || clock_timestamp()::text,
            'sha256'
        ),
        'hex'
    )
);

-- ALTER TABLE users DROP COLUMN api_key;