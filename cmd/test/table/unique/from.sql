CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1` int,
    `int2` int,
    `int3` int,
    `int4` int,
    `int5` int,
    `int6` int,
    # remained
    UNIQUE (`int1`),
    # modified
    UNIQUE KEY (`int2`, `int3`),
    # removed
    UNIQUE INDEX idx1 (`int4`),
    # renamed
    UNIQUE INDEX idx2 (`int6`)
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
    CONSTRAINT c1 UNIQUE KEY (`int1`),
    # modified
    CONSTRAINT c2 UNIQUE INDEX (`int2`, `int3`),
    # removed
    CONSTRAINT c3 UNIQUE INDEX idx1 (`int4`),
    # renamed
    CONSTRAINT c4 UNIQUE INDEX idx2 (`int6`)
);