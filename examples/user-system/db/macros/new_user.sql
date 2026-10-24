START TRANSACTION;

INSERT INTO Users (Username, PasswordHash, Salt)
VALUES (@Username, @Password, @Salt);

SELECT LAST_INSERT_ID() AS UserId;

COMMIT;