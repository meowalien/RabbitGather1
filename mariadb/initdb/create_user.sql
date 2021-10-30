create database `core` CHARACTER SET = utf8mb4;
CREATE USER 'core'@'%' IDENTIFIED BY '5678';
GRANT ALL ON core.* TO 'core'@'%';
FLUSH PRIVILEGES;

