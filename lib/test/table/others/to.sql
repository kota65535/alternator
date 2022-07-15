CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1`     int PRIMARY KEY AUTO_INCREMENT,
    `varchar1` varchar(20) UNIQUE,
    `varchar2` varchar(10)
);

CREATE TABLE `t2`
(
    # modified
    `varchar1` varchar(10),
    # renamed
    `varchar5` varchar(30),
    # moved
    `int1`     int,
    # added
    `varchar3` varchar(10),
    PRIMARY KEY (`int1`)
);
