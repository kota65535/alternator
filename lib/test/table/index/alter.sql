ALTER TABLE `db1`.`t1` CHANGE COLUMN `int7` `int8` int;
ALTER TABLE `db1`.`t1` DROP INDEX `idx1`;
ALTER TABLE `db1`.`t1` ADD INDEX `idx1` (`int5`);
ALTER TABLE `db1`.`t1` RENAME INDEX `idx2` TO `idx3`;
ALTER TABLE `db1`.`t1` ALTER INDEX `<unknown index name of 'INDEX (`int2`, `int3`)'>` INVISIBLE;