CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1` int,
    `int2` int,
    `int3` int,
    # to be retained
    CHECK (`int1` > 0),
    # to be modified
    CHECK (`int2` > 0),
    # to be dropped
    CHECK (`int3` > 0)
);

CREATE TABLE `t2`
(
    `int1` int,
    `int2` int,
    `int3` int,
    # to be retained
    CONSTRAINT c1 CHECK (`int1` > 0),
    # to be modified
    CONSTRAINT c2 CHECK (`int2` > 0),
    # to be dropped
    CONSTRAINT c3 CHECK (`int3` > 0)
);