create table t1
(
    bit1       bit,
    bit2       BIT(1)                         NOT NULL DEFAULT 1,
    tinyint1   tinyint,
    tinyint2   TINYINT(1) UNSIGNED ZEROFILL   NOT NULL DEFAULT 1,
    bool1      bool,
    bool2      BOOLEAN                        NOT NULL DEFAULT 1,
    smallint1  smallint,
    smallint2  smallINT(1) UNSIGNED ZEROFILL  NOT NULL DEFAULT 1,
    mediumint1 mediumint,
    mediumint2 MEDIUMINT(1) UNSIGNED ZEROFILL NOT NULL DEFAULT 1,
    `int1`     int,
    `int2`     INT(1) UNSIGNED ZEROFILL       NOT NULL DEFAULT 1,
    bigint1    bigint,
    bigint2    BIGINT(1) UNSIGNED ZEROFILL    NOT NULL DEFAULT 1,
    decimal1   decimal,
    decimal2   DECIMAL(2),
    decimal3   DEC(2, 1) UNSIGNED ZEROFILL    NOT NULL DEFAULT 1,
    float1     float,
    float2     FLOAT(2, 1) UNSIGNED ZEROFILL  NOT NULL DEFAULT 1.1,
    double1    double,
    double2    DOUBLE(2, 1) UNSIGNED ZEROFILL NOT NULL DEFAULT 1.1
);
