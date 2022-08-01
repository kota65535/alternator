CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1`     int AUTO_INCREMENT,
    `varchar1` varchar(20),
    `varchar2` varchar(10),
    PRIMARY KEY (`int1`)
)
    # larger value
    AUTO_INCREMENT = 100
