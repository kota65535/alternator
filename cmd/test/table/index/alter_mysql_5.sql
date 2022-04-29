ALTER TABLE `db1`.`t1` DROP INDEX `idx1`;
ALTER TABLE `db1`.`t1` ADD INDEX `idx1` (`int5`);
ALTER TABLE `db1`.`t1` RENAME INDEX `idx2` TO `idx3`;