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
	SET @Id_string = CONVERT(varchar(36), @Id);

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
	END CATCH
END
