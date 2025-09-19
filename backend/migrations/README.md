# migrations
Some notes on Postgres and our database tables.

## Resetting local postgres
In the beginning we're messing with the migration numbers a bit and so you'll get dirty version errors until we finalize the first few db tables. To reset your Postgres image, ensure there's no data you care about there, then:
```bash
make db-down # can't delete it while it's running
docker volume rm pin_pg_data
```

## Useful queries
Using the comments tree structure to find replies:
```sql
-- all replies to a comment
SELECT * FROM comments WHERE path <@ 'comment-1';
-- direct replies only
SELECT * FROM comments WHERE path ~ 'comment-1.*{1}';
```

User queries:
```sql
-- find user from google
SELECT u.* FROM users u 
JOIN oauth_accounts oa ON u.id = oa.user_id 
WHERE oa.provider = 'google' AND oa.provider_id = 'asdf';

-- is user logged in?
SELECT * FROM sessions WHERE session_token = 'asdf' AND expires_at > NOW();
```