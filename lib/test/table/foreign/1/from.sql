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
    FOREIGN KEY (`int2`) REFERENCES `t0` (`int1`),
    # removed
    FOREIGN KEY fk1 (`int3`) REFERENCES `t0` (`int1`),
    # recreated
    FOREIGN KEY fk2 (`int4`) REFERENCES `t0` (`int1`)
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
    CONSTRAINT c2 FOREIGN KEY (`int2`) REFERENCES `t0` (`int1`),
    # removed
    CONSTRAINT c3 FOREIGN KEY fk1 (`int3`) REFERENCES `t0` (`int1`),
    # recreated
    CONSTRAINT c4 FOREIGN KEY fk2 (`int4`) REFERENCES `t0` (`int1`)
);