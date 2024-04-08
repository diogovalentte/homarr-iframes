package cinemark

import "testing"

func TestGetInTheatersMovies(t *testing.T) {
	t.Run("should return a list of movies in theaters for a specific city", func(t *testing.T) {
		scraper := Scraper{}
		movies, err := scraper.GetInTheatersMovies("sao-paulo", -1, "716%2C690%2C699")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(movies) == 0 {
			t.Errorf("expected movies, got none")
		}

		println(len(movies))
		for _, movie := range movies {
			if movie.Title == "" {
				t.Errorf("expected movie title, got none")
			}

			if movie.CoverImgURL == "" {
				t.Errorf("expected movie cover, got none")
			}

		}
	})

	t.Run("should return a list of movies in theaters for a specific city with limit 5", func(t *testing.T) {
		scraper := Scraper{}
		movies, err := scraper.GetInTheatersMovies("sao-paulo", 5, "716")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(movies) == 0 {
			t.Errorf("expected movies, got none")
		}
		if len(movies) > 5 {
			t.Errorf("expected 5 movies, got %v", len(movies))
		}

		println(len(movies))
		for _, movie := range movies {
			if movie.Title == "" {
				t.Errorf("expected movie title, got none")
			}

			if movie.CoverImgURL == "" {
				t.Errorf("expected movie cover, got none")
			}

		}
	})

	t.Run("should return a list of movies in theaters for a specific city with limit 5 in all theaters", func(t *testing.T) {
		scraper := Scraper{}
		movies, err := scraper.GetInTheatersMovies("sao-paulo", 5, "")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(movies) == 0 {
			t.Errorf("expected movies, got none")
		}
		if len(movies) > 5 {
			t.Errorf("expected 5 movies, got %v", len(movies))
		}

		println(len(movies))
		for _, movie := range movies {
			if movie.Title == "" {
				t.Errorf("expected movie title, got none")
			}

			if movie.CoverImgURL == "" {
				t.Errorf("expected movie cover, got none")
			}

		}
	})
}
