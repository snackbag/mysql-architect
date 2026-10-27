INSERT INTO Users (Username, Password, Salt)
VALUES (@Username, @Password, @Salt);

SELECT LAST_INSERT_ID();