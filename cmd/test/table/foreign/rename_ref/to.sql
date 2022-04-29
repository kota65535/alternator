CREATE DATABASE db1;

USE db1;

CREATE TABLE `t0`
(
    # renamed
    `int10` int,
    # renamed
    `int20` int,
    PRIMARY KEY (`int10`),
    UNIQUE KEY (`int20`)
);

CREATE TABLE `t1`
(
    `int1` int,
    `int2` int,
    `int3` int,
    `int4` int,
    # retained
    FOREIGN KEY (`int1`) REFERENCES `t0` (`int10`),
    # modified
    FOREIGN KEY (`int2`) REFERENCES `t0` (`int20`) ON UPDATE CASCADE,
    # added
    FOREIGN KEY (`int4`) REFERENCES `t0` (`int20`)
);
