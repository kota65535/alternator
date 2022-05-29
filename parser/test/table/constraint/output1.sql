CREATE TABLE `t2`
(
    `int1` int,
    `int2` int,
    PRIMARY KEY (`int1` DESC),
    UNIQUE KEY (`int2` DESC),
    FOREIGN KEY (`int1`) REFERENCES `t1` (`int2`)
);