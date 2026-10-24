# MySQL Architect

Library to streamline and simplify MySQL database lifetime management and query building.

Databases still have to be created manually.

```mysql
CREATE DATABASE dbName;
CREATE USER 'username'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON dbName.* TO 'username'@'localhost';
```