# retained
CREATE DATABASE db1;

USE db1;

# renamed
CREATE TABLE `t13`
(
    `int1`     int,
    `varchar1` varchar(10) UNIQUE,
    PRIMARY KEY (`int1`)
);

# modified
CREATE TABLE `t12`
(
    `int1`     int,
    `int2`     int,
    # added
    `varchar1` varchar(10),
    PRIMARY KEY (`int1`),
    # retained, referencing the renamed table
    FOREIGN KEY (`int2`) REFERENCES `t13` (`int1`),
    # added, referencing the renamed table
    FOREIGN KEY (`varchar1`) REFERENCES `t13` (`varchar1`)
);

# modified
CREATE DATABASE `db2`
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_bin;

USE db2;

# retained
CREATE TABLE `t21`
(
    `int1`     int,
    # modified due to the change of the default collation
    `varchar1` varchar(10) UNIQUE,
    PRIMARY KEY (`int1`)
);

# modified
CREATE TABLE `t22`
(
    `int1`     int,
    `int2`     int,
    `varchar1` varchar(10),
    PRIMARY KEY (`int1`),
    # retained, referencing renamed table
    FOREIGN KEY (`int2`) REFERENCES `t21` (`int1`),
    # added, referencing renamed table
    FOREIGN KEY (`varchar1`) REFERENCES `t21` (`varchar1`)
);

# added
CREATE DATABASE `db3`
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_bin;

# added
CREATE TABLE `db3`.`t32`
(
    `varchar1` varchar(32),
    PRIMARY KEY (`varchar1`)
);

# added
CREATE TABLE `db3`.`t31`
(
    `int1` int,
    PRIMARY KEY (`int1`)
);

# db4 dropped
