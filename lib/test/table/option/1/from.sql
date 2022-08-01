CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `int1`     int AUTO_INCREMENT,
    `varchar1` varchar(20),
    `varchar2` varchar(10),
    PRIMARY KEY (`int1`)
)
    # remained
    KEY_BLOCK_SIZE = 1
    # modified
    MAX_ROWS = 100
    # removed (but remained)
    COMMENT = 'foo'



