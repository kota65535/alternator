CREATE TABLE `t2`
(
    `int1` int,
    `int2` int,
    PRIMARY KEY (`int1`),
    UNIQUE KEY (`int2`),
    FOREIGN KEY (`int1`) REFERENCES `t1` (`int2`)
);