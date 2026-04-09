IF EXISTS (
SELECT *
	FROM INFORMATION_SCHEMA.ROUTINES
WHERE SPECIFIC_SCHEMA = N'dbo'
	AND SPECIFIC_NAME = N'r_movie_by_id'
)
DROP PROCEDURE dbo.r_movie_by_id
GO
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
EXECUTE dbo.r_movie_by_id '98268a96-a6ac-444f-852a-c6472129aa22'
GO
