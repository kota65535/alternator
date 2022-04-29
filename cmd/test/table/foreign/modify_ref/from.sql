CREATE DATABASE db1;

USE db1;

CREATE TABLE `t0`
(
    # to be modified
    `t0_int1` int,
    # to be modified
    `t0_int2` int,
    PRIMARY KEY (`t0_int1`),
    UNIQUE KEY (`t0_int2`)
);

CREATE TABLE `t1`
(
    `t1_int1` int,
    `t1_int2` int,
    `t1_int3` int,
    `t1_int4` int,
    `t1_int5` int,
    # to be recreated
    CONSTRAINT `t1_ibfk_1` FOREIGN KEY (`t1_int1`) REFERENCES `t0` (`t0_int1`),
    # to be recreated
    CONSTRAINT `t1_ibfk_2` FOREIGN KEY (`t1_int2`) REFERENCES `t0` (`t0_int1`),
    # to be dropped
    CONSTRAINT `t1_ibfk_3` FOREIGN KEY fk1 (`t1_int3`) REFERENCES `t0` (`t0_int1`)
);
