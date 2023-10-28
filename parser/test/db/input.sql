create database db1;

# comment
-- comment
/* comment */

/****
  comment
*/
CREATE DATABASE IF NOT EXISTS `db2`
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_bin
    DEFAULT ENCRYPTION 'Y';

CREATE DATABASE IF NOT EXISTS `db3` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */
