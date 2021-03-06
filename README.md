# Alternator

Declarative SQL database schema management by human-friendly pure SQL.

- 100% pure SQL schema file
- Diff visualization between a schema file and actual database schemas
- Automatic generation and execution of ALTER statements to apply the schema
- No need to "import" redundant schemas from the database: write your own way

## Supported databases

- MySQL (8.x, 5.x)

## Install

### Mac

```bash
brew tap kota65535/alternator
brew install alternator
```

### Linux & Windows

Download the binary from [GitHub Releases](https://github.com/kota65535/alternator/releases/latest) and drop it in
your `$PATH`.

## Getting Started

1. Create a schema file. As an example here, we create `schema.sql` with the following content.

```sql
CREATE DATABASE example;
USE example;
CREATE TABLE users
(
    id   int PRIMARY KEY,
    name varchar(100)
);
CREATE TABLE blog_posts
(
    id        int PRIMARY KEY,
    title     varchar(100),
    body      text,
    author_id int,
    FOREIGN KEY (author_id) REFERENCES users (id)
);
```

2. Run the schema file to create a database and tables.

```sh
mysql -u root < schema.sql
```

3. Run `alternator plan` to verify the schema is up-to-date.

```sh
alternator plan schema.sql mysql://root@localhost/example
```

<details>
  <summary>Show Output</summary>

```
Fetching schemas of database 'example'...
――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
Schema diff:

  CREATE DATABASE `example`;

  CREATE TABLE `example`.`users`
  (
      `id`   int          NOT NULL,
      `name` varchar(100),
      PRIMARY KEY (`id`)
  );

  CREATE TABLE `example`.`blog_posts`
  (
      `id`        int          NOT NULL,
      `title`     varchar(100),
      `body`      text,
      `author_id` int,
      PRIMARY KEY (`id`),
      CONSTRAINT `blog_posts_ibfk_1` FOREIGN KEY `author_id` (`author_id`) REFERENCES `users` (`id`)
  );

――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
Your database schema is up-to-date! No change required.
```

</details>

4. Now let's modify the schema file.

```sql
CREATE DATABASE example;
USE example;
CREATE TABLE users
(
    id   int PRIMARY KEY AUTO_INCREMENT, # add AUTO_INCREMENT
    name varchar(100),
    INDEX (name)                         # added
);
CREATE TABLE blog_posts
(
    id        int PRIMARY KEY,
    title     varchar(200), # modify length 100 -> 200
    content   text,         # rename body -> content
    author_id int,
    FOREIGN KEY (author_id) REFERENCES users (id)
);
CREATE TABLE categories # added
(
    id   int PRIMARY KEY,
    name varchar(100)
);
```

5. Run `alternator plan` to show the schema diff and SQL statements that should be executed.
   Note that the foreign key `author_id` should be recreated because of the modification of the referencing key.

```sh
alternator plan schema.sql mysql://root@localhost/example
```

<details>
  <summary>Show Output</summary>

```diff
Fetching schemas of database 'example'...
――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
Schema diff:

  CREATE DATABASE `example`;

  CREATE TABLE `example`.`users`
  (
!     `id`   int          NOT NULL -> `id` int NOT NULL AUTO_INCREMENT,
      `name` varchar(100),
      PRIMARY KEY (`id`),
+     INDEX (`name`)
  );

  CREATE TABLE `example`.`blog_posts`
  (
      `id`        int          NOT NULL,
!     `title`     varchar(100)          -> `title`   varchar(200),
!     `body`      text                  -> `content` text,
      `author_id` int,
      PRIMARY KEY (`id`),
-     CONSTRAINT `blog_posts_ibfk_1` FOREIGN KEY `author_id` (`author_id`) REFERENCES `users` (`id`),
+     FOREIGN KEY (author_id) REFERENCES users (id)
  );

+ CREATE TABLE `example`.`categories`
+ (
+     `id`   int          NOT NULL,
+     `name` varchar(100),
+     PRIMARY KEY (`id`)
+ );

――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
Statements to execute:

ALTER TABLE `example`.`users` ADD INDEX (`name`);
ALTER TABLE `example`.`blog_posts` MODIFY COLUMN `title` varchar(200);
ALTER TABLE `example`.`blog_posts` CHANGE COLUMN `body` `content` text;
ALTER TABLE `example`.`blog_posts` DROP FOREIGN KEY `blog_posts_ibfk_1`;
ALTER TABLE `example`.`blog_posts` DROP INDEX `author_id`;
ALTER TABLE `example`.`users` MODIFY COLUMN `id` int NOT NULL AUTO_INCREMENT;
ALTER TABLE `example`.`blog_posts` ADD FOREIGN KEY (`author_id`) REFERENCES `users` (`id`);
CREATE TABLE `example`.`categories`
(
    `id`   int          NOT NULL,
    `name` varchar(100),
    PRIMARY KEY (`id`)
);
```

</details>

6. Run `alternator apply` to apply the schema change by executing planned SQL statements.

```sh
alternator apply schema.sql mysql://root@localhost/example
```

<details>
  <summary>Show Output</summary>

```
Fetching schemas of database 'example'...
――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
Statements to execute:

ALTER TABLE `example`.`users` ADD INDEX (`name`);
ALTER TABLE `example`.`blog_posts` MODIFY COLUMN `title` varchar(200);
ALTER TABLE `example`.`blog_posts` CHANGE COLUMN `body` `content` text;
ALTER TABLE `example`.`blog_posts` DROP FOREIGN KEY `blog_posts_ibfk_1`;
ALTER TABLE `example`.`blog_posts` DROP INDEX `author_id`;
ALTER TABLE `example`.`users` MODIFY COLUMN `id` int NOT NULL AUTO_INCREMENT;
ALTER TABLE `example`.`blog_posts` ADD FOREIGN KEY (`author_id`) REFERENCES `users` (`id`);
CREATE TABLE `example`.`categories`
(
    `id`   int          NOT NULL,
    `name` varchar(100),
    PRIMARY KEY (`id`)
);

――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
Do you want to apply? [y/n]: y
Executing: ALTER TABLE `example`.`users` ADD INDEX (`name`);
Executing: ALTER TABLE `example`.`blog_posts` MODIFY COLUMN `title` varchar(200);
Executing: ALTER TABLE `example`.`blog_posts` CHANGE COLUMN `body` `content` text;
Executing: ALTER TABLE `example`.`blog_posts` DROP FOREIGN KEY `blog_posts_ibfk_1`;
Executing: ALTER TABLE `example`.`blog_posts` DROP INDEX `author_id`;
Executing: ALTER TABLE `example`.`users` MODIFY COLUMN `id` int NOT NULL AUTO_INCREMENT;
Executing: ALTER TABLE `example`.`blog_posts` ADD FOREIGN KEY (`author_id`) REFERENCES `users` (`id`);
Executing: CREATE TABLE `example`.`categories`
(
    `id`   int          NOT NULL,
    `name` varchar(100),
    PRIMARY KEY (`id`)
);

Finished!
```

</details>
