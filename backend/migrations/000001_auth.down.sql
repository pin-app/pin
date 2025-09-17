DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
DROP TRIGGER IF EXISTS update_oauth_accounts_updated_at ON oauth_accounts;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS oauth_states;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS oauth_provider;