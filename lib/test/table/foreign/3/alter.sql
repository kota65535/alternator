ALTER TABLE `db1`.`t1` CHANGE COLUMN `int1` `int10` int;
ALTER TABLE `db1`.`t1` CHANGE COLUMN `int2` `int20` int;
ALTER TABLE `db1`.`t1` DROP FOREIGN KEY `t1_ibfk_2`;
ALTER TABLE `db1`.`t1` DROP INDEX `t1_ibfk_2`;
ALTER TABLE `db1`.`t1` ADD FOREIGN KEY (`int20`) REFERENCES `t0` (`int2`) ON UPDATE CASCADE;
ALTER TABLE `db1`.`t1` DROP FOREIGN KEY `t1_ibfk_3`;
ALTER TABLE `db1`.`t1` DROP INDEX `t1_ibfk_3`;
ALTER TABLE `db1`.`t1` ADD FOREIGN KEY `fk2` (`int5`) REFERENCES `t0` (`int2`);