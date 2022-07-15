CREATE TABLE `t2`
(
    `int1`  int,
    `date1` date
)
    PARTITION BY RANGE (year(`date1`))
    PARTITIONS 3
    SUBPARTITION BY HASH (to_days(`date1`))
    SUBPARTITIONS 2
    (
        PARTITION p0 VALUES LESS THAN (1990),
        PARTITION p1 VALUES LESS THAN (2000),
        PARTITION p2 VALUES LESS THAN (MAXVALUE)
    );