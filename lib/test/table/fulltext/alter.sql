ALTER TABLE `db1`.`t1` CHANGE COLUMN `var7` `var8` varchar(16);
ALTER TABLE `db1`.`t1` DROP INDEX `idx1`;
ALTER TABLE `db1`.`t1` ADD FULLTEXT INDEX `idx1` (`var5`);
ALTER TABLE `db1`.`t1` RENAME INDEX `idx2` TO `idx3`;
ALTER TABLE `db1`.`t1` ALTER INDEX `<unknown index name of 'FULLTEXT INDEX (`var2`, `var3`)'>` INVISIBLE;