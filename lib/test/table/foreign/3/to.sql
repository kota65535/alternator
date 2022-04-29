CREATE DATABASE db1;

USE db1;

CREATE TABLE `t0`
(
    `int1` int,
    `int2` int,
    PRIMARY KEY (`int1`),
    UNIQUE KEY (`int2`)
);

CREATE TABLE `t1`
(
    `int10` int,
    `int20` int,
    `int3` int,
    `int4` int,
    `int5` int,
    `int6` int,
    # retained
    FOREIGN KEY (`int10`) REFERENCES `t0` (`int1`),
    # modified
    FOREIGN KEY (`int20`) REFERENCES `t0` (`int2`) ON UPDATE CASCADE,
    # added
    FOREIGN KEY fk2 (`int5`) REFERENCES `t0` (`int2`)
);
