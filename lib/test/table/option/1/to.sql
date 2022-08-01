CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1`     int PRIMARY KEY AUTO_INCREMENT,
    `varchar1` varchar(20),
    `varchar2` varchar(10)
)
    # added
    AUTO_INCREMENT = 10
    # modified
    MAX_ROWS = 200
