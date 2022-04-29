CREATE DATABASE db1;

USE db1;

CREATE TABLE `t0`
(
    # modified
    `t0_int1` int AUTO_INCREMENT,
    # modified
    `t0_int2` int NOT NULL,
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
    # recreated
    FOREIGN KEY (`t1_int1`) REFERENCES `t0` (`t0_int1`),
    # recreated
    FOREIGN KEY (`t1_int2`) REFERENCES `t0` (`t0_int2`) ON UPDATE CASCADE,
    # added
    FOREIGN KEY (`t1_int5`) REFERENCES `t0` (`t0_int2`)
);
