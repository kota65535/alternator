# to be retained charset & collate
CREATE DATABASE db1
    DEFAULT CHARSET = utf8mb4
    DEFAULT COLLATE = utf8mb4_bin;

USE db1;

# to be retained charset & collate
CREATE TABLE `t11`
(
    `int1`     int,
    # to be retained
    `varchar1` varchar(10),
    # to be added charset, collate
    `varchar2` varchar(10),
    # to be modified charset, collate
    `varchar3` varchar(10) CHARSET sjis COLLATE sjis_japanese_ci,
    # to be deleted charset, collate
    `varchar4` varchar(10) CHARSET sjis COLLATE sjis_japanese_ci,

    PRIMARY KEY (`int1`)
);

# to be added charset & collate
CREATE DATABASE db2;

USE db2;

# to be added charset & collate
CREATE TABLE `t21`
(
    `int1`     int,
    # to be retained
    `varchar1` varchar(10),
    # to be added charset, collate
    `varchar2` varchar(10),
    # to be modified charset, collate
    `varchar3` varchar(10) CHARSET sjis COLLATE sjis_japanese_ci,
    # to be deleted charset, collate
    `varchar4` varchar(10) CHARSET sjis COLLATE sjis_japanese_ci,

    PRIMARY KEY (`int1`)
);

# to be modified charset & collate
CREATE DATABASE db3
    DEFAULT CHARSET = utf16
    DEFAULT COLLATE = utf16_general_ci;

USE db3;

# to be modified charset & collate
CREATE TABLE `t31`
(
    `int1`     int,
    # to be retained
    `varchar1` varchar(10),
    # to be added charset, collate
    `varchar2` varchar(10),
    # to be modified charset, collate
    `varchar3` varchar(10) CHARSET sjis COLLATE sjis_japanese_ci,
    # to be deleted charset, collate
    `varchar4` varchar(10) CHARSET sjis COLLATE sjis_japanese_ci,

    PRIMARY KEY (`int1`)
);
