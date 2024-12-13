CREATE TABLE `news` (
  `id` varchar(255) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `autor` varchar(255) DEFAULT NULL,
  `text` varchar(255) DEFAULT NULL,
  `publishdate` datetime DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1

CREATE TABLE news 
(
id INT NOT NULL AUTO_INCREMENT,
title VARCHAR(255)  NOT NULL ,
autor VARCHAR(255)  NOT NULL ,
content TEXT,
publishdate DATETIME NOT NULL,
 PRIMARY KEY (id)
)