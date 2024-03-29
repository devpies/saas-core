CREATE TABLE IF NOT EXISTS customers (
    customer_id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    payment_method VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE INDEX idx_customer_surname ON customers(last_name);
CREATE INDEX idx_customer_email ON customers(email);

CREATE TABLE IF NOT EXISTS subscriptions (
    subscription_id VARCHAR(36) PRIMARY KEY,
    plan VARCHAR(36) NOT NULL,
    transaction_id VARCHAR(36) NOT NULL,
    subscription_status_id int NOT NULL,
    amount INT NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY(customer_id) REFERENCES customers (customer_id)
);

CREATE INDEX idx_subscription_transaction ON subscriptions(transaction_id);
CREATE INDEX idx_subscription_status_id ON subscriptions(subscription_status_id);
CREATE INDEX idx_subscription_customer ON subscriptions(tenant_id);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id VARCHAR(36) PRIMARY KEY,
    amount INT NOT NULL,
    currency VARCHAR(36) NOT NULL,
    last_four VARCHAR(4) NOT NULL,
    bank_return_code VARCHAR(255) NULL,
    transaction_status_id int NOT NULL,
    expiration_month INT NOT NULL,
    expiration_year INT NOT NULL,
    subscription_id VARCHAR(255) NOT NULL,
    payment_intent VARCHAR(255) NULL,
    payment_method VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY(subscription_id) REFERENCES subscriptions (subscription_id)
);

CREATE INDEX idx_transaction_status_id ON transactions(transaction_status_id);

CREATE TABLE IF NOT EXISTS transaction_statuses (
    transaction_status_id int PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE IF NOT EXISTS subscription_statuses (
    subscription_status_id int PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

-- enable RLS
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;

--create policies
CREATE POLICY transactions_isolation_policy ON transactions
    USING (tenant_id = current_setting('app.current_tenant'));
CREATE POLICY subscriptions_isolation_policy ON subscriptions
    USING (tenant_id = current_setting('app.current_tenant'));
CREATE POLICY customers_isolation_policy ON customers
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE USER user_a WITH PASSWORD 'postgres';
GRANT ALL ON ALL TABLES IN SCHEMA "public" TO user_a;
GRANT postgres to user_a;