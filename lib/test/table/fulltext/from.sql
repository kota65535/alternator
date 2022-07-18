CREATE DATABASE db1;

USE db1;

CREATE TABLE `t1`
(
    `var1` varchar(16),
    `var2` varchar(16),
    `var3` varchar(16),
    `var4` varchar(16),
    `var5` varchar(16),
    `var6` varchar(16),
    # remained
    FULLTEXT KEY (`var1`),
    # modified
    FULLTEXT INDEX (`var2`, `var3`),
    # removed
    FULLTEXT INDEX idx1 (`var4`),
    # renamed
    FULLTEXT INDEX idx2 (`var6`)
);
