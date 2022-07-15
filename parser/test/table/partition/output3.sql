CREATE TABLE `t3`
(
    `double1` double,
    `double2` double
)
    PARTITION BY RANGE COLUMNS (`double1`, `double2`)
    PARTITIONS 3
    (
        PARTITION p0 VALUES LESS THAN (1990.1),
        PARTITION p1 VALUES LESS THAN (2000.1),
        PARTITION p2 VALUES LESS THAN (MAXVALUE)
    );