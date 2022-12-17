DROP TABLE IF EXISTS users;
source users.sql;

INSERT INTO users(userName, fullName, email, hashedPassword)
VALUES ('test', 'Test User', 'test@email', '$2a$10$2bLycFqUmc6m6iLkaeUgKOGwzekGd9IoAPMbXRNNuJ8Sv9ItgV29O');

INSERT INTO users(userName, fullName, email, hashedPassword, admin)
VALUES ('admin', 'Admin User', 'admin@email', '$2a$10$2bLycFqUmc6m6iLkaeUgKOGwzekGd9IoAPMbXRNNuJ8Sv9ItgV29O', 1);

DROP TABLE IF EXISTS tokens;
source tokens.sql;
