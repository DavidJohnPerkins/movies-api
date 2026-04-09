IF EXISTS (
SELECT *
	FROM INFORMATION_SCHEMA.ROUTINES
WHERE SPECIFIC_SCHEMA = N'dbo'
	AND SPECIFIC_NAME = N'r_movie'
)
DROP PROCEDURE dbo.r_movie
GO
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
EXECUTE dbo.r_movie
GO
