package main

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/ifan/webdev/views"
)

func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}

var projects []views.ProjectData = []views.ProjectData{
	{"Project 1", "/projects/1", "static/img/placeholder.svg?height=200&width=200", "This is a description of a project."},
	{"Project 2", "/projects/2", "static/img/placeholder.svg?height=200&width=200", "This is a description of a project."},
	{"Project 3", "/projects/3", "static/img/placeholder.svg?height=200&width=200", "This is a description of a project."},
	{"Project 4", "/projects/4", "static/img/placeholder.svg?height=200&width=200", "This is a description of a project."},
	{"Project 5", "/projects/5", "static/img/placeholder.svg?height=200&width=200", "This is a description of a project."},
}

func getIndex(c *gin.Context) {
	render(c, 200, views.Home(projects))
}

func getProject(c *gin.Context) {
	name := c.Param("name")
	i, err := strconv.Atoi(name)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid project number"})
		return
	}
	if i < 1 || i > len(projects) {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}
	render(c, 200, views.ProjectPage(projects[i-1]))
}

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", getIndex)
	r.GET("/projects/:name", getProject)
	r.Run("localhost:8080")
}
