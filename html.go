package main

import (
	"html/template"
	"io"
	"os"
)

func copyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Use io.Copy to copy the contents from source to destination
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Flush file to disk
	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

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

	// add css
	copyFile("./styles.css", "./public/styles.css")

	// Confirm successful file generation
	println("Static HTML file generated successfully.")
}
