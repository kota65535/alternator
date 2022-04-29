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
    # renamed
    `int10` int,
    # renamed
    `int20` int,
    # renamed
    `int30` int,
    # renamed
    `int40` int,
    # retained
    FOREIGN KEY (`int10`) REFERENCES `t0` (`int1`),
    # modified
    FOREIGN KEY (`int20`) REFERENCES `t0` (`int2`) ON UPDATE CASCADE,
    # added
    FOREIGN KEY (`int40`) REFERENCES `t0` (`int2`)
);
