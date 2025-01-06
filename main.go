package main

import (
	"strconv"

	"database/sql"
	"log"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/ifan/webdev/views"
	_ "github.com/mattn/go-sqlite3"
)

func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}

var db *sql.DB

type Project struct {
	Name    string
	Image   string
	Content string
}

func seedDatabase() {
	projects := []Project{
		{Name: "Project 1", Image: "static/img/placeholder.svg?height=200&width=200", Content: "This is a description of a project."},
		{Name: "Project 2", Image: "static/img/placeholder.svg?height=200&width=200", Content: "This is a description of a project."},
		{Name: "Project 3", Image: "static/img/placeholder.svg?height=200&width=200", Content: "This is a description of a project."},
		{Name: "Project 4", Image: "static/img/placeholder.svg?height=200&width=200", Content: "This is a description of a project."},
		{Name: "Project 5", Image: "static/img/placeholder.svg?height=200&width=200", Content: "This is a description of a project."},
	}

	for _, project := range projects {
		_, err := db.Exec("INSERT INTO projects (name, image, content) VALUES (?, ?, ?)", project.Name, project.Image, project.Content)
		if err != nil {
			log.Printf("Failed to insert project %s: %v", project.Name, err)
		}
	}
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./projects.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		image TEXT,
		content TEXT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		seedDatabase()
	}
}

func getProjects() ([]views.ProjectData, error) {
	rows, err := db.Query("SELECT id, name, image, content FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []views.ProjectData
	for rows.Next() {
		var project views.ProjectData
		if err := rows.Scan(&project.Id, &project.Name, &project.Image, &project.Content); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func getProjectByID(id int) (views.ProjectData, error) {
	var project views.ProjectData
	err := db.QueryRow("SELECT id, name, image, content FROM projects WHERE id = ?", id).Scan(&project.Id, &project.Name, &project.Image, &project.Content)
	if err != nil {
		return project, err
	}
	return project, nil
}

func getIndex(c *gin.Context) {
	projects, err := getProjects()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get projects. Error: " + err.Error()})
		render(c, 500, views.Error("Failed to get projects. Error: "+err.Error()))
	}
	render(c, 200, views.Home(projects))
}

func getProject(c *gin.Context) {
	name := c.Param("id")
	id, err := strconv.Atoi(name)
	if err != nil {
		render(c, 500, views.Error("Failed to convert project ID to integer. Error: "+err.Error()))
		return
	}

	project, err := getProjectByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			render(c, 404, views.NotFound())
		} else {
			render(c, 500, views.Error("Failed to get project. Error: "+err.Error()))
		}
		return
	}
	render(c, 200, views.ProjectPage(project))
}

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", getIndex)
	r.GET("/projects/:id", getProject)
	r.GET("/create", createProject)
	r.POST("/create", saveProject)
	r.Run("localhost:8080")
}

func createProject(c *gin.Context) {
	render(c, 200, views.CreateProjectPage())
}

func saveProject(c *gin.Context) {
	name := c.PostForm("name")
	image := c.PostForm("image")
	content := c.PostForm("content")

	_, err := db.Exec("INSERT INTO projects (name, image, content) VALUES (?, ?, ?)", name, image, content)
	if err != nil {
		render(c, 500, views.Error("Failed to save project. Error: "+err.Error()))
		return
	}

	render(c, 200, views.Success("Project saved successfully."))
}
