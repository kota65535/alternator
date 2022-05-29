CREATE DATABASE db1
    DEFAULT CHARSET utf8mb4;


USE db1;

CREATE TABLE `t1`
(
    `int1` int,
    `int2` int,
    `int3` int,
    `int4` int,
    `int5` int,
    `var6` varchar(10),
    # remained
    KEY (`int1`),
    # modified
    INDEX (`int2`, `int3`),
    # removed
    INDEX idx1 (`int4`),
    # renamed
    INDEX idx2 (`var6`(5))
);
