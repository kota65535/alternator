ALTER TABLE `db1`.`t2` MODIFY COLUMN `varchar1` varchar(10);
ALTER TABLE `db1`.`t2` CHANGE COLUMN `varchar4` `varchar5` varchar(30);
ALTER TABLE `db1`.`t2` DROP COLUMN `varchar2`;
ALTER TABLE `db1`.`t2` MODIFY COLUMN `int1` int NOT NULL AFTER `varchar5`;
ALTER TABLE `db1`.`t2` MODIFY COLUMN `int2` int DEFAULT '1' AFTER `int1`;
ALTER TABLE `db1`.`t2` ADD COLUMN `varchar3` varchar(10) AFTER `int2`;