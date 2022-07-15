CREATE DATABASE db1;

USE db1;

# column options are converted to each create definition, so no change is expected
CREATE TABLE `t1`
(
    `int1`     int AUTO_INCREMENT,
    `json1`    json NOT NULL,
    `geo1`     geometry GENERATED ALWAYS
    PRIMARY KEY (`int1`)
);
