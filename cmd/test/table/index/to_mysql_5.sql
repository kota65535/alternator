CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1` int,
    `int2` int,
    `int3` int,
    `int4` int,
    `int5` int,
    `var6` varchar(10),
    # modified
    INDEX (`int2`, `int3`),
    # remained
    KEY (`int1`),
    # added
    INDEX idx1 (`int5`),
    # renamed
    INDEX idx3 (`var6`(5) ASC)
);
