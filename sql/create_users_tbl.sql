DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `userName` varchar(30) NOT NULL,
  `fullName` varchar(70) NOT NULL,
  `email` varchar(256) NOT NULL,
  `hashedPassword` char(72) NOT NULL,
  `resetToken` varchar(44) DEFAULT NULL,
  `resetExpires` datetime DEFAULT NULL,
  PRIMARY KEY (`userName`),
  UNIQUE KEY `email` (`email`)
);

INSERT INTO users(userName, fullName, email, hashedPassword)
VALUES ("test", "go test user", "test@email", "$2a$10$2bLycFqUmc6m6iLkaeUgKOGwzekGd9IoAPMbXRNNuJ8Sv9ItgV29O")
