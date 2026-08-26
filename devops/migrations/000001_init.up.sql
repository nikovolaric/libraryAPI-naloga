CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT NOT NULL,
    last_name  TEXT NOT NULL
);

CREATE TABLE books (
    id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title     TEXT NOT NULL,
    total     INTEGER NOT NULL DEFAULT 0 CHECK (total >= 0),
    available INTEGER NOT NULL DEFAULT 0 CHECK (available >= 0 AND available <= total)
);

CREATE TABLE loans (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id     UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    borrow_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    return_date TIMESTAMPTZ
);

CREATE INDEX idx_loans_user_id ON loans(user_id);
CREATE INDEX idx_loans_book_id ON loans(book_id);
