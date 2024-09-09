package main

import (
	"html/template"
	"os"
)

func genHtml() {
	data := struct {
		RefreshTimeMs uint32
	}{
		RefreshTimeMs: REFRESH_TIME*1000 + 60000,
	}

	tmpl, err := os.ReadFile("./template.html")
	if err != nil {
		panic(err)
	}

	// Create a new template and parse the defined HTML
	t := template.Must(template.New("staticPage").Parse(string(tmpl)))

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
