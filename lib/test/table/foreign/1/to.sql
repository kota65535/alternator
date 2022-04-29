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
    # remained
    FOREIGN KEY (`int1`) REFERENCES `t0` (`int1`),
    # recreated
    FOREIGN KEY (`int2`) REFERENCES `t0` (`int1`) ON UPDATE CASCADE,
    # added
    FOREIGN KEY fk1 (`int5`) REFERENCES `t0` (`int2`),
    # recreated
    FOREIGN KEY fk3 (`int6`)  REFERENCES `t0` (`int1`)
);

CREATE TABLE `t2`
(
    `int1` int,
    `int2` int,
    `int3` int,
    `int4` int,
    `int5` int,
    `int6` int,
    # remained
    CONSTRAINT c1 FOREIGN KEY (`int1`) REFERENCES `t0` (`int1`),
    # recreated
    CONSTRAINT c2 FOREIGN KEY (`int2`) REFERENCES `t0` (`int2`) ON UPDATE CASCADE,
    # added
    CONSTRAINT c3 FOREIGN KEY fk1 (`int5`) REFERENCES `t0` (`int2`),
    # recreated
    CONSTRAINT c4 FOREIGN KEY fk3 (`int6`)  REFERENCES `t0` (`int1`)
);