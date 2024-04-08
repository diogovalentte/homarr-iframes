package cinemark

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

// GetInTheatersMovies returns in theater movies (current movies playing on the theater) for a specific city
// If the cinemark site doesn't have a page for the city, it will return the page of São Paulo instead of
// 404, so it's better to check if the city is valid before calling this function
func (s *Scraper) GetInTheatersMovies(city string, limit int, theaters string) ([]Movie, error) {
	var sharedErr error
	var movies []Movie
	var err error
	nextPage := true

	url := fmt.Sprintf("%s/%s/filmes/em-cartaz", siteURL, city)
	if theaters != "" {
		url = fmt.Sprintf("%s?cinema=%s", url, theaters)
	}
	count := 0
	for nextPage {
		s.resetCollector()
		nextPage = false

		s.c.OnHTML("section.movies div.active > div.row > div", func(e *colly.HTMLElement) {
			if limit != -1 && count >= limit {
				return
			}
			movie := Movie{}

			movie.Title = e.ChildAttr("article > div > a", "title")
			movie.Title = strings.TrimSpace(strings.Replace(movie.Title, "Filme", "", -1))
			movie.URL = e.ChildAttr("article > div > a", "href")
			movie.URL = fmt.Sprintf("%s%s", siteURL, movie.URL)
			movie.CoverImgURL = e.ChildAttr("article > div > a > picture > source", "srcset")
			movie.AgeRating = e.ChildText("article > div > div.movie-details > div.movie-rating > span.rating-abbr")
			movie.AgeRatingColor = getMovieAgeRatingColor(movie.AgeRating)

			movieLabel := e.ChildText("article > div > span.movie-label")
			if movieLabel != "" {
				movie.Label = movieLabel
			}

			movies = append(movies, movie)
			count++
		})

		s.c.OnHTML("section.movies div.active > nav > ul > li > a.pagination-next", func(e *colly.HTMLElement) {
			if limit != -1 && count >= limit {
				return
			}
			nextPage = true
			url = e.Attr("href")
		})

		s.c.Visit(url)
		if err != nil {
			return nil, err
		}
	}

	return movies, sharedErr
}
