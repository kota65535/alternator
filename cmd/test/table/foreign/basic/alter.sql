ALTER TABLE `db1`.`t2` DROP FOREIGN KEY `c2`;
ALTER TABLE `db1`.`t2` DROP INDEX `c2`;
ALTER TABLE `db1`.`t2` ADD CONSTRAINT `c2` FOREIGN KEY `c2` (`int2`) REFERENCES `t0` (`int2`) ON UPDATE CASCADE;
ALTER TABLE `db1`.`t2` DROP FOREIGN KEY `c3`;
ALTER TABLE `db1`.`t2` DROP INDEX `c3`;
ALTER TABLE `db1`.`t2` DROP FOREIGN KEY `c4`;
ALTER TABLE `db1`.`t2` DROP INDEX `c4`;
ALTER TABLE `db1`.`t2` ADD CONSTRAINT `c3` FOREIGN KEY `c3` (`int5`) REFERENCES `t0` (`int2`);
ALTER TABLE `db1`.`t2` ADD CONSTRAINT `c4` FOREIGN KEY `c4` (`int6`) REFERENCES `t0` (`int1`);