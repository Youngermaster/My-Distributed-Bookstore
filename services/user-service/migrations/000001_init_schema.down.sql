-- Drop tables in reverse order of creation (respecting foreign key constraints)
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF NOT EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;

-- Drop UUID extension (optional, comment out if other services use it)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
