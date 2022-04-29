ALTER TABLE `db1`.`t1` DROP INDEX `idx1`;
ALTER TABLE `db1`.`t1` RENAME INDEX `idx2` TO `idx3`;
ALTER TABLE `db1`.`t1` ALTER INDEX `<unknown index name of 'UNIQUE KEY (int2, int3)'>` INVISIBLE;
ALTER TABLE `db1`.`t2` DROP INDEX `idx1`;
ALTER TABLE `db1`.`t2` RENAME INDEX `idx2` TO `idx3`;
ALTER TABLE `db1`.`t2` ALTER INDEX `c2` INVISIBLE;
ALTER TABLE `db1`.`t1` ADD UNIQUE KEY `idx2` (`int5`);
ALTER TABLE `db1`.`t2` ADD UNIQUE KEY `idx2` (`int5`);