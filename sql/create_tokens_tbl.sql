DROP TABLE IF EXISTS `tokens`;

CREATE TABLE `tokens` (
  `value` varchar(44) NOT NULL,
  `expires` datetime NOT NULL,
  `type` varchar(7) NOT NULL,
  `userName` varchar(30) NOT NULL,
  PRIMARY KEY (`value`)
);
