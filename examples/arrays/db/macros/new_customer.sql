INSERT INTO Customers (Name, RegistrationDate)
VALUES (@Name, @RegistrationDate);

SELECT LAST_INSERT_ID();