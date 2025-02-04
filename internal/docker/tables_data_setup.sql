-- Create the books table
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title TEXT UNIQUE NOT NULL,
    available_copies INT NOT NULL CHECK (available_copies >= 0)
);

-- Insert initial book records
INSERT INTO books (title, available_copies) VALUES
    ('book1', 5),
    ('book2', 3),
    ('book3', 1),
    ('book4', 0)
ON CONFLICT (title) DO NOTHING; -- Prevent duplicate inserts

CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    borrower_name TEXT NOT NULL,
    loan_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    return_date TIMESTAMP NOT NULL,
    is_returned BOOLEAN DEFAULT FALSE
);

CREATE UNIQUE INDEX unique_active_loan
    ON loans (book_id, borrower_name)
    WHERE is_returned = FALSE;

