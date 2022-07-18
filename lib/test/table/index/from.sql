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
    # remained
    KEY (`int1`),
    # modified
    INDEX (`int2`, `int3`),
    # removed
    INDEX idx1 (`int4`),
    # renamed
    INDEX idx2 ((`int6`*2))
);
