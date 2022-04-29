CREATE TABLE `t1`
(
    `int1`     int,
    `int2`     int,
    `int3`     int,
    `varchar1` varchar(10),
    `varchar2` varchar(10),
    `varchar3` varchar(10),
    INDEX (`int1`),
    INDEX `idx1` (`int2`, `int3`),
    INDEX `idx2` (`varchar1`) USING BTREE KEY_BLOCK_SIZE 1 COMMENT 'foo' VISIBLE,
    FULLTEXT INDEX (`varchar2`),
    FULLTEXT INDEX `idx3` (`varchar3`) KEY_BLOCK_SIZE 1 WITH PARSER ngram COMMENT 'foo' VISIBLE
);