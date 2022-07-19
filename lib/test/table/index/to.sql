CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1` int,
    `int2` int,
    `int3` int,
    `int4` int,
    `int5` int,
    `int6` int,
    `int8` int,
    # remained
    KEY (`int1`),
    # modified
    INDEX (`int2`, `int3`) INVISIBLE,
    # added
    INDEX idx1 (`int5`),
    # renamed
    INDEX idx3 ((`int6` * 2)),
    # column renamed
    INDEX ((`int8` * 3))
);
