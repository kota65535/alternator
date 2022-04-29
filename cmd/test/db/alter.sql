ALTER TABLE `db2`.`t22` ADD COLUMN `varchar1` varchar(10) AFTER `int2`;
ALTER TABLE `db2`.`t21` RENAME TO `db2`.`t23`;
ALTER TABLE `db2`.`t22` ADD FOREIGN KEY (`varchar1`) REFERENCES `t23` (`varchar1`);
CREATE DATABASE `db3`
    DEFAULT CHARACTER SET = utf8mb4
    DEFAULT COLLATE = utf8mb4_bin;
CREATE TABLE `db3`.`t32`
(
    `varchar1` varchar(32) NOT NULL,
    PRIMARY KEY (`varchar1`)
);
CREATE TABLE `db3`.`t31`
(
    `int1` int NOT NULL,
    PRIMARY KEY (`int1`)
);
DROP DATABASE `db4`;