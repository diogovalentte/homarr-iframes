package cinemark

import "testing"

func TestCinemark_OnDisplayByTheater(t *testing.T) {
	theaterIds := []int{715, 710} // Cinemark Shopping Eldorado
	limit := 100
	limitProvided := true
	t.Run("should return a list of movies in a specific theater", func(t *testing.T) {
		c := Cinemark{}
		movies, err := c.GetOnDisplayByTheater(theaterIds, limit, limitProvided)
		if err != nil {
			t.Fatalf("error getting movies: %v", err)
		}

		if len(movies) == 0 {
			t.Fatalf("no movies found")
		}
	})
}
