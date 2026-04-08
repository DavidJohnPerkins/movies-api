-- Drop the stored procedure if it already exists
IF EXISTS (
SELECT *
	FROM INFORMATION_SCHEMA.ROUTINES
WHERE SPECIFIC_SCHEMA = N'dbo'
	AND SPECIFIC_NAME = N'r_movie_by_id'
)
DROP PROCEDURE dbo.r_movie_by_id
GO
-- Create the stored procedure in the specified schema
CREATE PROCEDURE dbo.r_movie_by_id
	@id UNIQUEIDENTIFIER
AS
	SELECT
		m.Id, 
		m.Title, 
		m.Director,
		m.ReleaseDate,
		m.TicketPrice,
		m.CreatedAt,
		m.UpdatedAt
	FROM
		dbo.movies m
	WHERE 
		m.Id = @id
GO
-- example to execute the stored procedure we just created
EXECUTE dbo.r_movie_by_id '98268a96-a6ac-444f-852a-c6472129aa22'
--EXECUTE dbo.r_movies ''
GO

-- Drop the stored procedure if it already exists
IF EXISTS (
SELECT *
	FROM INFORMATION_SCHEMA.ROUTINES
WHERE SPECIFIC_SCHEMA = N'dbo'
	AND SPECIFIC_NAME = N'r_movie'
)
DROP PROCEDURE dbo.r_movie
GO
-- Create the stored procedure in the specified schema
CREATE PROCEDURE dbo.r_movie
AS
	SELECT
		m.Id, 
		m.Title, 
		m.Director,
		m.ReleaseDate,
		m.TicketPrice,
		m.CreatedAt,
		m.UpdatedAt
	FROM
		dbo.movies m
	ORDER BY 
		m.Id
GO
-- example to execute the stored procedure we just created
EXECUTE dbo.r_movie
--EXECUTE dbo.r_movies ''
GO

IF EXISTS (
SELECT *
	FROM INFORMATION_SCHEMA.ROUTINES
WHERE SPECIFIC_SCHEMA = N'dbo'
	AND SPECIFIC_NAME = N'c_movie'
)
DROP PROCEDURE dbo.c_movie
GO
CREATE PROCEDURE dbo.c_movie
    @json NVARCHAR(MAX)
AS
BEGIN
    DECLARE @Id UNIQUEIDENTIFIER,
    		@Id_string varchar(36)

    SELECT @Id = Id
    FROM OPENJSON(@json)
    WITH (
        Id UNIQUEIDENTIFIER '$.id'
    );

	SET @Id_string = CONVERT(varchar(36), @Id)

	BEGIN TRY
		-- Now you can use @Id however you want
		INSERT INTO dbo.Movies (Id, Title, Director, ReleaseDate, TicketPrice, CreatedAt, UpdatedAt)
		SELECT *,
		GETDATE(),
		GETDATE()
		FROM OPENJSON(@json)
		WITH (
			Id UNIQUEIDENTIFIER '$.id',
			Title NVARCHAR(200) '$.title',
			Director NVARCHAR(200) '$.director',
			ReleaseDate DATETIME2 '$.release_date',
			TicketPrice DECIMAL(10,2) '$.ticket_price'
		);
	END TRY
	BEGIN CATCH
  		RAISERROR ('Duplicate key: %s - operation failed', 16, 1, @Id_string)
--	END CATCH
END
EXEC c_movie '{ "id": "98268a96-a6ac-444f-852a-c6472129aa30", "title": "Most of the Time", "director": "David Perkins", "release_date": "2021-07-01T01:01:01.00Z", "ticket_price": 100.00 }'

