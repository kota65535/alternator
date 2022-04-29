create table t1
(
    `int1`   INT,
    `int2`   INT,
    `int3`   INT,
    varchar1 VARCHAR(10),
    varchar2 VARCHAR(10),
    varchar3 VARCHAR(10),

    INDEX (`int1`),
    INDEX idx1 (`int2`, `int3`),
    INDEX idx2 (varchar1)
        USING BTREE
        KEY_BLOCK_SIZE = 1
        VISIBLE
        COMMENT 'foo',
    FULLTEXT INDEX (varchar2),
    FULLTEXT INDEX idx3 (varchar3) WITH PARSER ngram KEY_BLOCK_SIZE = 1 VISIBLE COMMENT 'foo'
);
