CREATE DATABASE db1;

USE db1;

# column options are converted to each create definition, so no change is expected
CREATE TABLE `t1`
(
    `int1`     int AUTO_INCREMENT,
    `varchar1` varchar(20),
    `varchar2` varchar(10) AS (concat(`varchar1`,'foo')),
    PRIMARY KEY (`int1`),
    UNIQUE KEY (`varchar1`)
);

CREATE TABLE `t2`
(
    # to be moved
    `int1`     int,
    # to be retained
    `int2`     int DEFAULT 1,
    # to be modified
    `varchar1` varchar(20),
    # to be renamed
    `varchar4` varchar(30),
    # to be dropped
    `varchar2` varchar(10),
    PRIMARY KEY (`int1`)
);
