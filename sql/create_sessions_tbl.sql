DROP TABLE IF EXISTS `sessions`;

CREATE TABLE `sessions` (
  `token` varchar(44) NOT NULL,
  `expires` datetime NOT NULL,
  `userName` varchar(30) NOT NULL,
  PRIMARY KEY (`token`)
);
