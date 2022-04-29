CREATE DATABASE db1;

USE db1;

CREATE TABLE `t3`
(
    `int1` int,
    `int2` int,
    PRIMARY KEY (`int1`),
    UNIQUE KEY (`int2`)
);

CREATE TABLE `t1`
(
    `int1` int,
    `int2` int,
    `int3` int,
    `int4` int,
    `int5` int,
    `int6` int,
    # retained
    FOREIGN KEY (`int1`) REFERENCES `t3` (`int1`),
    # modified
    FOREIGN KEY (`int2`) REFERENCES `t3` (`int2`) ON UPDATE CASCADE,
    # added
    FOREIGN KEY fk1 (`int5`) REFERENCES `t3` (`int2`)
);
