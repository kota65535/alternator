create table t1
(
    char1       char,
    char2       CHARACTER(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin    NOT NULL DEFAULT 'a',
    varchar1    varchar(1),
    varchar2    VARCHAR(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin      NOT NULL DEFAULT 'a',
    binary1     binary,
    binary2     BINARY(1)                                                 NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    varbinary1  varbinary(1),
    varbinary2  VARBINARY(2)                                              NOT NULL DEFAULT 'a',
    tinyblob1   tinyblob,
    tinyblob2   TINYBLOB                                                  NOT NULL,
    tinytext1   tinytext,
    tinytext2   TINYTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin        NOT NULL,
    blob1       blob,
    blob2       BLOB(1)                                                   NOT NULL,
    text1       text,
    text2       TEXT(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin         NOT NULL,
    mediumblob1 mediumblob,
    mediumblob2 MEDIUMBLOB                                                NOT NULL,
    mediumtext1 mediumtext,
    mediumtext2 MEDIUMTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin      NOT NULL,
    longblob1   longblob,
    longblob2   LONGBLOB                                                  NOT NULL,
    longtext1   longtext,
    longtext2   LONGTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin        NOT NULL,
    enum1       enum ('a'),
    enum2       ENUM ('a', 'b') CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT 'a',
    set1        set ('a'),
    set2        SET ('a', 'b') CHARACTER SET utf8mb4 COLLATE utf8mb4_bin  NOT NULL DEFAULT 'a'
);

