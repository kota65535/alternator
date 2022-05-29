create table t1
(
    `int1` INT,
    `int2` INT,
    INDEX (`int2`)
);

create table t2
(
    `int1` INT,
    `int2` INT,
    PRIMARY KEY (`int1` DESC),
    UNIQUE KEY (`int2` DESC),
    FOREIGN KEY (`int1`) REFERENCES t1 (`int2`)
);

create table t3
(
    `int1` INT,
    `int2` INT,
    CONSTRAINT u1 PRIMARY KEY (`int1`, `int2`)
        USING BTREE
        KEY_BLOCK_SIZE = 1
        COMMENT 'foo'
        VISIBLE,
    CONSTRAINT u2 UNIQUE KEY (`int1` ASC, `int2`)
        USING BTREE
        KEY_BLOCK_SIZE = 1
        COMMENT 'foo'
        VISIBLE,
    CONSTRAINT u3 FOREIGN KEY i3 (`int1` ASC, `int2`) REFERENCES t1 (`int1`, `int2`)
        MATCH FULL
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
