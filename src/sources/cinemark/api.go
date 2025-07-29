package cinemark

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	apiURL             = "https://br-www-frontend-ext-prod.cinemark.com.br/bff-api/v1/"
	defaultMoviesLimit = 999
)

func (c *Cinemark) GetOnDisplayByTheater(theaterIDs []int, limit int, limitProvided bool) ([]Movie, error) {
	if !limitProvided {
		limit = defaultMoviesLimit
	} else {
		if limit < 1 {
			return []Movie{}, nil
		}
	}

	moviesNames := make(map[string]struct{})
	moviesSlice := []Movie{}
	for _, theaterID := range theaterIDs {
		var responseData onDisplayByTheaterResponse
		err := c.baseRequest("GET", apiURL+"movies/OnDisplayByTheater"+fmt.Sprintf("?&theaterId=%d&pageNumber=1&pageSize=%d", theaterID, limit), nil, &responseData)
		if err != nil {
			return nil, err
		}

		for _, movie := range responseData.Movies {
			if _, exists := moviesNames[movie.Name]; !exists {
				movieURL := fmt.Sprintf("https://www.cinemark.com.br/filme/%s", movie.Slug)
				if len(theaterIDs) > 1 {
					movieURL = movieURL + "?city=true"
				}

				moviesSlice = append(moviesSlice, Movie{
					Name:           movie.Name,
					CoverImgURL:    movie.Assets[0].URL,
					URL:            movieURL,
					AgeRating:      movie.AgeIndication,
					AgeRatingColor: getMovieAgeRatingColor(movie.AgeIndication),
					Genre:          movie.Genre,
				})
				moviesNames[movie.Name] = struct{}{}
			}
		}
	}

	return moviesSlice, nil
}

type onDisplayByTheaterResponse struct {
	Movies []struct {
		Name          string `json:"name"`
		Slug          string `json:"slug"`
		AgeIndication string `json:"ageIndication"`
		Genre         string `json:"genre"`
		Assets        []struct {
			Type int    `json:"type"`
			URL  string `json:"url"`
		} `json:"assets"`
	} `json:"dataResult"`
}

func (Cinemark) baseRequest(method, url string, body io.Reader, target any) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.cinemark.com.br")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status (%s): %s", resp.Status, string(resBody))
	}

	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(resBody, target); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\nreponse text: %s", err.Error(), string(resBody))
	}

	return nil
}

func getMovieAgeRatingColor(rating string) string {
	rating = strings.ToUpper(rating)
	switch rating {
	case "L":
		return "#00bb22"
	case "10", "A10":
		return "#5891cd"
	case "12", "A12":
		return "#edcb0c"
	case "14", "A14":
		return "#f6962d"
	case "16", "A16":
		return "#dd021c"
	case "18", "A18":
		return "#000"
	default:
		return "gray"
	}
}
