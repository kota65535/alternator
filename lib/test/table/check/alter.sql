ALTER TABLE `db1`.`t1` ALTER CHECK `<unknown constraint name of '`int2` > 0'>` NOT ENFORCED;
ALTER TABLE `db1`.`t1` DROP CHECK `<unknown constraint name of (`int3` > 0)>`;
ALTER TABLE `db1`.`t1` ADD CHECK (`int3` > 2);
ALTER TABLE `db1`.`t2` ALTER CHECK `c2` NOT ENFORCED;
ALTER TABLE `db1`.`t2` DROP CHECK `c3`;
ALTER TABLE `db1`.`t2` ADD CONSTRAINT `c3` CHECK (`int3` > 2);