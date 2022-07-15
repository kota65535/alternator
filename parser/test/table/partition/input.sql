create table t1
(
    `int1` int,
    primary key (`int1`)
)
    partition by hash (`int1`);

create table t2
(
    `int1`  int,
    `date1` date
)
    partition by range (year(`date1`))
        partitions 3
        subpartition by hash (to_days(`date1`))
        subpartitions 2
        (
        partition p0 values less than (1990),
        partition p1 values less than (2000),
        partition p2 values less than (maxvalue)
        );

create table t3
(
    `double1` double,
    `double2` double
)
    partition by range columns (`double1`, `double2`)
        partitions 3
        (
        partition p0 values less than (1990.1),
        partition p1 values less than (2000.1),
        partition p2 values less than maxvalue
        );
