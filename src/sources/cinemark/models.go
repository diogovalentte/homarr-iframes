package cinemark

// Movie is the struct for a movie
type Movie struct {
	// Name is the movie name
	Name string
	// CoverImgURL is the URL for the movie cover image
	CoverImgURL string
	// URL is the URL for the movie details page
	URL string
	// AgeRating is the movie rating, can be: L, 12, 14, 16, 18
	AgeRating string
	// AgeRatingColor is the color for the movie rating
	AgeRatingColor string
	// Genre is the movie genre
	Genre string
}
