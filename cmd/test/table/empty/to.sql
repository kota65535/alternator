CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1`     int,
    `varchar1` varchar(20),
    `varchar2` varchar(10) DEFAULT 'foo',
    PRIMARY KEY (`int1`)
);

CREATE TABLE `t2`
(
    `varchar1` varchar(10),
    `varchar4` varchar(30),
    `int1`     int,
    `varchar3` varchar(10),
    PRIMARY KEY (`int1`)
);
