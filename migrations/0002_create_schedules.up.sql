CREATE TABLE schedules (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    interval INT,
    days_of_month INT[],
    specific_dates DATE[],
    even_odd TEXT,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);