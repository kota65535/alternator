ALTER TABLE `db1`.`t0` MODIFY COLUMN `t0_int2` int NOT NULL;
ALTER TABLE `db1`.`t1` DROP FOREIGN KEY `t1_ibfk_1`;
ALTER TABLE `db1`.`t1` DROP INDEX `t1_ibfk_1`;
ALTER TABLE `db1`.`t1` DROP FOREIGN KEY `t1_ibfk_2`;
ALTER TABLE `db1`.`t1` DROP INDEX `t1_ibfk_2`;
ALTER TABLE `db1`.`t1` ADD FOREIGN KEY (`t1_int2`) REFERENCES `t0` (`t0_int2`) ON UPDATE CASCADE;
ALTER TABLE `db1`.`t1` DROP FOREIGN KEY `t1_ibfk_3`;
ALTER TABLE `db1`.`t1` DROP INDEX `t1_ibfk_3`;
ALTER TABLE `db1`.`t0` MODIFY COLUMN `t0_int1` int NOT NULL AUTO_INCREMENT;
ALTER TABLE `db1`.`t1` ADD FOREIGN KEY (`t1_int1`) REFERENCES `t0` (`t0_int1`);
ALTER TABLE `db1`.`t1` ADD FOREIGN KEY `fk2` (`t1_int5`) REFERENCES `t0` (`t0_int2`);