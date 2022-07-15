CREATE TABLE `t1`
(
    `int1` int,
    PRIMARY KEY (`int1`)
)
    PARTITION BY HASH (`int1`);