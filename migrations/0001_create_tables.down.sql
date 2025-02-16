
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_transactions_sender;
DROP INDEX IF EXISTS idx_transactions_receiver;
DROP INDEX IF EXISTS idx_inventory_user;


DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS store;
DROP TABLE IF EXISTS users;
