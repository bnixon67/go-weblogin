DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `userName` varchar(30) NOT NULL,
  `fullName` varchar(70) NOT NULL,
  `email` varchar(256) NOT NULL,
  `hashedPassword` char(72) NOT NULL,
  `sessionToken` varchar(44) DEFAULT NULL,
  `sessionExpires` datetime DEFAULT NULL,
  `resetToken` varchar(44) DEFAULT NULL,
  `resetExpires` datetime DEFAULT NULL,
  PRIMARY KEY (`userName`),
  UNIQUE KEY `email` (`email`)
);
