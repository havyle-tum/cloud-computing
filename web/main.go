package main

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Defines a "model" that we can use to communicate with the
// frontend or the database
type Book struct {
	MongoID     primitive.ObjectID `bson:"_id,omitempty"`
	ID          string
	BookName    string
	BookAuthor  string
	BookEdition string
	BookPages   string
	BookYear    string
}

// Wraps the "Template" struct to associate a necessary method
// to determine the rendering procedure
type Template struct {
	tmpl *template.Template
}

// Preload the available templates for the view folder.
// This builds a local "database" of all available "blocks"
// to render upon request, i.e., replace the respective
// variable or expression.
// For more on templating, visit https://jinja.palletsprojects.com/en/3.0.x/templates/
// to get to know more about templating
// You can also read Golang's documentation on their templating
// https://pkg.go.dev/text/template
func loadTemplates() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

// Method definition of the required "Render" to be passed for the Rendering
// engine.
// Contraire to method declaration, such syntax defines methods for a given
// struct. "Interfaces" and "structs" can have methods associated with it.
// The difference lies that interfaces declare methods whether struct only
// implement them, i.e., only define them. Such differentiation is important
// for a compiler to ensure types provide implementations of such methods.
func (t *Template) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

// fetchBooksFromAPI calls the books microservice to get all books
func fetchBooksFromAPI() ([]Book) {	
	resp, err := http.Get("http://books-get:8081/api/books")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var books []Book
	err = json.Unmarshal(body, &books)
	if err != nil {
		return nil
	}

	return books
}

func fetchAuthorsFromAPI() ([]string) {	
	resp, err := http.Get("http://books-get:8081/api/authors")
	if err != nil {
		return []string{} // or return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []string{} // or return nil
	}

	var authors []string
	err = json.NewDecoder(resp.Body).Decode(&authors)
	if err != nil {
		return []string{} // or return nil
	}

	return authors
}

func fetchYearsFromAPI() ([]string) {	
	resp, err := http.Get("http://books-get:8081/api/years")
	if err != nil {
		return []string{} // or return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []string{} // or return nil
	}

	var years []string
	err = json.NewDecoder(resp.Body).Decode(&years)
	if err != nil {
		return []string{} // or return nil
	}

	return years
}

func main() {
	// Here we prepare the server
	e := echo.New()

	// Define our custom renderer
	e.Renderer = loadTemplates()

	// Log the requests. Please have a look at echo's documentation on more
	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("/css", "css")

	// Endpoint definition. Here, we divided into two groups: top-level routes
	// starting with /, which usually serve webpages. For our RESTful endpoints,
	// we prefix the route with /api to indicate more information or resources
	// are available under such route.
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.GET("/books", func(c echo.Context) error {
		books := fetchBooksFromAPI()
		return c.Render(200, "book-table", books)
	})

	e.GET("/authors", func(c echo.Context) error {
		authors := fetchAuthorsFromAPI()
		return c.Render(200, "author-table", authors)
	})

	e.GET("/years", func(c echo.Context) error {
		years := fetchYearsFromAPI()
		return c.Render(200, "year-table", years)
	})

	e.GET("/search", func(c echo.Context) error {
		return c.Render(200, "search-bar", nil)
	})

	e.GET("/create", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	// We start the server and bind it to port 8080. For future references, this
	// is the application's port and not the external one. For this first exercise,
	// they could be the same if you use a Cloud Provider. If you use ngrok or similar,
	// they might differ.
	// In the submission website for this exercise, you will have to provide the internet-reachable
	// endpoint: http://<host>:<external-port>
	e.Logger.Fatal(e.Start(":8080"))
}