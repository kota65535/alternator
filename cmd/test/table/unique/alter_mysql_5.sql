ALTER TABLE `db1`.`t1` DROP INDEX `idx1`;
ALTER TABLE `db1`.`t1` RENAME INDEX `idx2` TO `idx3`;
ALTER TABLE `db1`.`t2` DROP INDEX `idx1`;
ALTER TABLE `db1`.`t2` RENAME INDEX `idx2` TO `idx3`;
ALTER TABLE `db1`.`t1` ADD UNIQUE KEY `idx2` (`int5`);
ALTER TABLE `db1`.`t2` ADD UNIQUE KEY `idx2` (`int5`);