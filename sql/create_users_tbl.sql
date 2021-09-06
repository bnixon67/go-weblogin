DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `userName` varchar(30) NOT NULL,
  `hashedPassword` char(72) NOT NULL,
  `sessionToken` varchar(36) DEFAULT NULL,
  `firstName` varchar(50) DEFAULT NULL,
  `lastName` varchar(50) DEFAULT NULL,
  `email` varchar(256) DEFAULT NULL,
  `sessionExpires` datetime DEFAULT NULL,
  PRIMARY KEY (`userName`)
);

GRANT SELECT, INSERT, UPDATE, DELETE ON `users` TO `user`@`%`;
