create table `t1`
(
    `int1` int,
    PRIMARY KEY (`int1`)
)
    PARTITION BY RANGE (`int1`) (
        PARTITION p0 VALUES LESS THAN (6),
        PARTITION p1 VALUES LESS THAN (11),
        PARTITION p2 VALUES LESS THAN (16),
        PARTITION p3 VALUES LESS THAN (21)
        );

