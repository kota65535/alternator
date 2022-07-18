# retained charset & collate
CREATE DATABASE db1
    DEFAULT CHARSET = utf8mb4
    DEFAULT COLLATE = utf8mb4_bin;

USE db1;

# retained charset & collate
CREATE TABLE `t11`
(
    `int1`     int,
    # retained
    `varchar1` varchar(10),
    # added charset, collate
    `varchar2` varchar(10) CHARSET sjis COLLATE sjis_japanese_ci,
    # modified charset, collate
    `varchar3` varchar(10) CHARSET utf16 COLLATE utf16_bin,
    # deleted charset, collate
    `varchar4` varchar(10),

    PRIMARY KEY (`int1`)
);

# added charset & collate
CREATE DATABASE db2
    DEFAULT CHARSET = utf8mb4
    DEFAULT COLLATE = utf8mb4_bin;

USE db2;

# added charset & collate
CREATE TABLE `t21`
(
    `int1`     int,
    # retained
    `varchar1` varchar(10),
    # added charset, collate
    `varchar2` varchar(10) CHARSET sjis COLLATE sjis_japanese_ci,
    # modified charset, collate
    `varchar3` varchar(10) CHARSET utf16 COLLATE utf16_bin,
    # deleted charset, collate
    `varchar4` varchar(10),

    PRIMARY KEY (`int1`)
);

# modified charset & collate
CREATE DATABASE db3
    DEFAULT CHARSET = utf8mb4
    DEFAULT COLLATE = utf8mb4_bin;

USE db3;

# modified charset & collate
CREATE TABLE `t31`
(
    `int1`     int,
    # retained
    `varchar1` varchar(10),
    # added charset, collate
    `varchar2` varchar(10) CHARSET sjis COLLATE sjis_japanese_ci,
    # modified charset, collate
    `varchar3` varchar(10) CHARSET utf16 COLLATE utf16_bin,
    # deleted charset, collate
    `varchar4` varchar(10),

    PRIMARY KEY (`int1`)
);
