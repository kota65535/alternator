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
    `int1` int,
    `int2` int,
    `int3` int,
    `int4` int,
    `int5` int,
    `int6` int,
    # to be retained
    CONSTRAINT `t1_ibfk_1` FOREIGN KEY (`int1`) REFERENCES `t0` (`int1`),
    # to be modified
    CONSTRAINT `t1_ibfk_2` FOREIGN KEY (`int2`) REFERENCES `t0` (`int1`),
    # to be dropped
    CONSTRAINT `t1_ibfk_3` FOREIGN KEY fk1 (`int3`) REFERENCES `t0` (`int1`)
);
