CREATE TABLE `t1`
(
    `date1`      date,
    `date2`      date         NOT NULL DEFAULT (CURRENT_DATE + INTERVAL 1 YEAR),
    `datetime1`  datetime,
    `datetime2`  datetime(1)  NOT NULL,
    `timestamp1` timestamp,
    `timestamp2` timestamp(1) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `time1`      time,
    `time2`      time(1)      NOT NULL,
    `year1`      year,
    `year2`      year(4)      NOT NULL
);