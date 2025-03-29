CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE loans (
    id SERIAL PRIMARY KEY,
    user_id INTEGER not null,
    amount DECIMAL(10, 2) not null,
    interest DECIMAL(5, 2) not null,
    week INTEGER not null,
    total DECIMAL(10, 2) not null,
    status INTEGER not null,
    weekly_payment DECIMAL(10, 2) not null,
    start_date TIMESTAMP not null,
    created_at TIMESTAMP not null,
	updated_at TIMESTAMP not null
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER not null,
    user_id INTEGER not null,
    week_number INTEGER not null,
    amount DECIMAL(10, 2) not null,
    status INTEGER not null,
    due_date TIMESTAMP not null,
    paid_date TIMESTAMP,
    created_at TIMESTAMP not null,
	updated_at TIMESTAMP not null
);