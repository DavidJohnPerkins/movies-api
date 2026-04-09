IF EXISTS (
SELECT *
	FROM INFORMATION_SCHEMA.ROUTINES
WHERE SPECIFIC_SCHEMA = N'dbo'
	AND SPECIFIC_NAME = N'd_movie'
)
DROP PROCEDURE dbo.u_movie
GO
CREATE PROCEDURE dbo.d_movie
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
        DELETE 
            m
		FROM
			dbo.Movies m
		WHERE
			m.Id = @Id
	END TRY
	BEGIN CATCH
  		RAISERROR ('Delete for key: %s - operation failed', 16, 1, @Id_string)
	END CATCH
END
