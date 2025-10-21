CREATE TABLE IF NOT EXISTS payments (
    id              VARCHAR NOT NULL,
    stripe_id        VARCHAR NOT NULL,
    amount          INT NOT NULL,
    currency        VARCHAR NOT NULL,
    status          VARCHAR NOT NULL,
    email           VARCHAR NOT NULL UNIQUE,
    payment_method   VARCHAR NOT NULL,
    idempotency_key  VARCHAR NOT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_payments_id PRIMARY KEY (id),
    CONSTRAINT uq_payments_email UNIQUE (email)
    );
