package sources

import "fmt"

// GetBaseNothingToShowiFrame returns an HTML code for when there is nothing to show
// The template is a background image with some message
func GetBaseNothingToShowiFrame(backgroundColor, backgroundImageURL, backgroundPosition, backgroundSize, brightness string) []byte {
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
	html = fmt.Sprintf(html, backgroundColor, backgroundImageURL, backgroundPosition, backgroundSize, brightness)

	return []byte(html)
}
