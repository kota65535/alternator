CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1`     int PRIMARY KEY AUTO_INCREMENT,
    `varchar1` varchar(20) UNIQUE,
    `varchar2` varchar(10) AS (concat(`varchar1`, 'foo'))
);

CREATE TABLE `t2`
(
    # modified
    `varchar1` varchar(10),
    # renamed
    `varchar5` varchar(30),
    # moved
    `int1`     int,
    # retained
    `int2`     int DEFAULT '1',
    # added
    `varchar3` varchar(10),
    PRIMARY KEY (`int1`)
);
