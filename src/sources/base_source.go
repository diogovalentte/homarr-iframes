package sources

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// GetBaseNothingToShowiFrame returns an HTML code for when there is nothing to show
// The template is a background image with some message
func GetBaseNothingToShowiFrame(backgroundColor, backgroundImageURL, backgroundPosition, backgroundSize, brightness, APIURLPath string) []byte {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>No Movies To Show</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            text-align: center;
            height: 100vh;
        }

        .background-image {
            background-color: %s;
            background-image: url('%s');
            background-position: %s;
            background-size: %s;
            position: absolute;
            filter: brightness(%s);
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            z-index: -1;
        }
    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = '%s';
                const response = await fetch(url);
                const data = await response.json();

                if (lastHash === null) {
                    lastHash = data.hash;
                } else {
                    if (data.hash !== lastHash) {
                        lastHash = data.hash;
                        location.reload();
                    }
                }
            } catch (error) {
                console.error('Error getting last update from the API:', error);
            }
        }

        function fetchAndUpdate() {
            fetchData();
            setTimeout(fetchAndUpdate, 10000); // 10 seconds
        }

        %s
    </script>
</head>
<body>
    <div class="background-image"></div>
</body>
</html>
    `
	if backgroundColor == "light" {
		backgroundColor = "#ffffff"
	} else if backgroundColor == "dark" {
		backgroundColor = "#25262b"
	}
	var refreshIframe string
	if APIURLPath != "" {
		refreshIframe = "fetchAndUpdate();"
	}
	html = fmt.Sprintf(html, backgroundColor, backgroundImageURL, backgroundPosition, backgroundSize, brightness, APIURLPath, refreshIframe)

	return []byte(html)
}

// GetHash returns a SHA256 hash of a object
func GetHash(s interface{}) [32]byte {
	structString := fmt.Sprintf("%v", s)
	return sha256.Sum256([]byte(structString))
}

// IsToday returns true if the date is today
func IsToday(date time.Time) bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(today) && date.Before(today.Add(24*time.Hour))
}
