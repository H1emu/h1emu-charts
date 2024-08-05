package main

import (
	"html/template"
	"os"
)

const tmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>H1emu Charts</title>
</head>
<body>
  <iframe src="./connections.html" width="100%" height="600px" frameborder="0">
        Your browser does not support iframes.
    </iframe>
  <iframe src="./Top Killer Main EU 1.html" width="100%" height="600px" frameborder="0">
        Your browser does not support iframes.
    </iframe>
  <iframe src="./Top Killer Main US 1.html" width="100%" height="600px" frameborder="0">
        Your browser does not support iframes.
    </iframe>
  <iframe src="./Top Zombie Killer Help.html" width="100%" height="600px" frameborder="0">
        Your browser does not support iframes.
    </iframe>
  
    {{.Connections}}

  <script> setTimeout(()=>location.reload(),{{.RefreshTimeMs}})</script>
</body>
</html>
`

func genHtml() {
	data := struct {
		Connections   string
		RefreshTimeMs uint32
	}{
		RefreshTimeMs: REFRESH_TIME*1000 + 60000,
	}

	// Create a new template and parse the defined HTML
	t := template.Must(template.New("staticPage").Parse(tmpl))

	// Create or open a file to save the generated HTML
	file, err := os.Create("public/index.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Execute the template and write the output to the file
	err = t.Execute(file, data)
	if err != nil {
		panic(err)
	}

	// Confirm successful file generation
	println("Static HTML file generated successfully.")
}
