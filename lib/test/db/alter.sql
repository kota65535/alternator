ALTER TABLE `db1`.`t11` RENAME TO `db1`.`t13`;
ALTER TABLE `db1`.`t12` ADD COLUMN `varchar1` varchar(10) AFTER `int2`;
ALTER TABLE `db1`.`t12` ADD FOREIGN KEY (`varchar1`) REFERENCES `t13` (`varchar1`);
ALTER DATABASE `db2` DEFAULT COLLATE = utf8mb4_bin;
ALTER TABLE `db2`.`t21` DEFAULT COLLATE = utf8mb4_bin;
ALTER TABLE `db2`.`t21` MODIFY COLUMN `varchar1` varchar(10);
ALTER TABLE `db2`.`t22` DEFAULT COLLATE = utf8mb4_bin;
ALTER TABLE `db2`.`t22` ADD COLUMN `varchar1` varchar(10) AFTER `int2`;
ALTER TABLE `db2`.`t22` ADD FOREIGN KEY (`varchar1`) REFERENCES `t21` (`varchar1`);
CREATE DATABASE `db3`
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