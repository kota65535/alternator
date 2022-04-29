create table t1
(
    date1      date,
    date2      DATE         NOT NULL,
    datetime1  datetime,
    datetime2  DATETIME(1)  NOT NULL,
    timestamp1 timestamp,
    timestamp2 TIMESTAMP(1) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    time1      time,
    time2      TIME(1)      NOT NULL,
    year1      year,
    year2      YEAR(4)      NOT NULL
);

