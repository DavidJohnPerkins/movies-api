IF EXISTS (
SELECT *
	FROM INFORMATION_SCHEMA.ROUTINES
WHERE SPECIFIC_SCHEMA = N'dbo'
	AND SPECIFIC_NAME = N'u_movie'
)
DROP PROCEDURE dbo.u_movie
GO
CREATE PROCEDURE dbo.u_movie
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
		WITH w_values AS (
			SELECT 
				Id,
				Title,
				Director,
				ReleaseDate,
				TicketPrice,
				GETDATE() AS UpdatedAt
			FROM OPENJSON(@json)
			WITH (
				Id UNIQUEIDENTIFIER '$.id',
				Title NVARCHAR(200) '$.title',
				Director NVARCHAR(200) '$.director',
				ReleaseDate DATETIME2 '$.release_date',
				TicketPrice DECIMAL(10,2) '$.ticket_price'
			)
		)
		UPDATE
			m
		SET
			m.Title = w.Title,
			m.Director = w.Director,
			m.ReleaseDate = w.ReleaseDate,
			m.TicketPrice = w.TicketPrice,
			m.UpdatedAt = w.UpdatedAt
		FROM
			dbo.Movies m,
			w_values w
		WHERE
			m.Id = @Id
	END TRY
	BEGIN CATCH
  		RAISERROR ('Update for key: %s - operation failed', 16, 1, @Id_string)
	END CATCH
END
