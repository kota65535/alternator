# to be retained
CREATE DATABASE db1;

USE db1;

# to be retained
CREATE TABLE `t11`
(
    `int1`     int,
    `varchar1` varchar(10),
    PRIMARY KEY (`int1`)
);

# to be dropped
CREATE DATABASE `db4`
    DEFAULT CHARACTER SET = utf8mb4
    DEFAULT COLLATE = utf8mb4_bin;

USE db4;

# to be dropped
CREATE TABLE `t41`
(
    `int1`     int,
    `varchar1` varchar(10),
    PRIMARY KEY (`int1`)
);

# to be modified
CREATE DATABASE `db2`
    DEFAULT CHARSET utf8mb4;

# to be renamed
CREATE TABLE `db2`.`t21`
(
    `int1`     int,
    `varchar1` varchar(10) UNIQUE,
    PRIMARY KEY (`int1`)
);

# to be modified
CREATE TABLE `db2`.`t22`
(
    `int1` int,
    `int2` int,
    PRIMARY KEY (`int1`),
    # to be retained, referencing table to rename
    FOREIGN KEY (`int2`) REFERENCES `t21` (`int1`)
);

# db3 will be added
