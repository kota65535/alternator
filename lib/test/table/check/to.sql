CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1` int,
    `int2` int,
    `int3` int,
    # remained
    CHECK (`int1` > 0),
    # recreated
    CHECK (`int2` > 0) NOT ENFORCED,
    # added
    CHECK (`int3` > 2)
);

CREATE TABLE `t2`
(
    `int1` int,
    `int2` int,
    `int3` int,
    # remained
    CONSTRAINT c1 CHECK (`int1` > 0),
    # recreated
    CONSTRAINT c2 CHECK (`int2` > 0) NOT ENFORCED,
    # added
    CONSTRAINT c3 CHECK (`int3` > 2)
);