package sources

import "fmt"

func GetBaseNoItemsToShowiFrame(itemName, textColor, backgroundImageURL, theme string) []byte {
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
        height: 100vh;
        background-color: %s;
    }

    .background-image {
        background-image: url('%s');
        background-position: center;
        background-size: cover;
        position: absolute;
        filter: brightness(0.3);
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: -1;
        border-radius: 10px;
    }

    .text {
        color: %s;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
        font-weight: bold;
        font-size: 3rem;
    }
</style>
</head>
<body>
    <div class="background-image"></div>
    <div class="text">No %s to show</div>
</body>
</html>
    `
	backgroundColor := "#ffffff"
	if theme == "dark" {
		backgroundColor = "#25262b"
	}
	html = fmt.Sprintf(html, backgroundColor, backgroundImageURL, textColor, itemName)

	return []byte(html)
}
