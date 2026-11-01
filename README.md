# MySQL Architect

Library to streamline and simplify MySQL database lifetime management and query building, by
abstracting complicated multi-query tasks into macros, which are easy to maintain and review.

Databases still have to be created manually.

```mysql
CREATE DATABASE dbName;
CREATE USER 'username'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON dbName.* TO 'username'@'localhost';
```

## Usage

For examples, see [/examples](/examples).

MySQL Architect manages tables for you, which are defined in your architecture. An architecture
is defined by an `architecture.json` file inside a directory. This file registers all tables and
macros.

You should write every path is with forward slashes, as they are translated automatically.

```
{
    "tables": ["tables/Users.json"],
    "macros": ["macros/get_user_by_name", "macros/get_user_by_id", "macros/new_user"]
}
```

You can define tables and macros in any directory, as long as they are inside the architecture
root.

```go
import architect "github.com/snackbag/mysql-architect"

func main() {
// load all macros and tables
arch, err := architect.LoadArchitecture("db")
if err != nil {
panic(err)
}

// open db connection and create all tables
db, err := arch.Init(username, password, dbName)
if err != nil {
panic(err)
}
}
```

### Tables

A table's JSON file is similar to its schema, and is built like so:

```json
{
    "columns": [
        {
            "name": "UserId",
            "type": "int",
            "primary": true,
            "autoIncrement": true
        }
    ]
}
```

Columns must have a `name` and `type`. Additionally, they may also be `primary` (bool), `unique`
(bool), `autoIncrement` (bool) and `nullable` (bool). You can create foreign keys by using
`references`, as shown in the
[/examples/basic-tables](/examples/basic-tables/db/tables/UserPermissions.json) example.

MySQL Architect supports text, varchar, all sizes of integer (+unsigned integers), booleans, date,
time, datetime and timestamp. More data can be given to types by using a JSON object instead of
a string, which supports `length` (for varchar) and `fsp` (for time related things). For example:

```json
{
    "name": "Username",
    "type": {
        "name": "varchar",
        "length": 255
    }
}
```

### Macros

Macros let you execute (multiple) SQL scripts inside one transaction. Every call of a macro
starts a transaction and commits it after all queries successfully execute.

```json
{
    "params": {
        "UserId": {
            "type": "int"
        },
        "Username": {
            "type": {
                "name": "varchar",
                "length": 255
            }
        }
    },
    "execute": [
        {
            "file": "macros/new_user.sql",
            "values": {
                "Username": "Username",
                "Password": "Password"
            },
            "overwriteParams": true
        },
        {
            "file": "macros/set_admin.sql",
            "values": {
                "UserId": "UserId"
            }
        }
    ]
}
```

```sql
INSERT INTO Users (Username, Password)
VALUES (@Username, @Password);

SELECT LAST_INSERT_ID() AS UserId;
```

```sql
INSERT INTO UserPermissions (UserId, Permission)
VALUES (@UserId, 'admin');
```

All macros must be registered in `architecture.json` as seen in the introduction. Macros can
take parameters, which can be assigned to values inside scripts. The `values` key in script
definitions takes a parameter name and assigns that name to the value in the script. When the
macro is executed, MySQL Architect ensures the values are inserted in the correct order.

```go
// execute macro
macro := arch.Macros["macros/new_admin"]
result, err := macro.Exec(db, map[string]any{"Username": "admin", "Password": "spinach"})
if err != nil {
panic(err)
}

// get user id
userId, _ = macros.Get[int](result, "UserId")
```

You can use `macros.Get` to get columns returned by the macro. When `overwriteParams` is enabled,
columns replace the parameters for the next script. That means that if the first script returns
`UserId` and a parameter named `UserId` exists, the parameter's value gets overwritten with the
new value.

For more information about multiple returned values,
see [/examples/basic-tables](/examples/basic-tables).