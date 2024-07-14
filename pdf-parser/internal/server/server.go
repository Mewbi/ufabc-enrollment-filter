package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"pdf-parser/config"
	"pdf-parser/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type controller struct {
	service *service.Service
}

func New(s *service.Service) *controller {
	return &controller{
		service: s,
	}
}

type InputEnrollment struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := config.Get()
		c.Header("Access-Control-Allow-Origin", config.Server.CorsHost)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, User-Agent",
		)
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (ct *controller) health(c *gin.Context) {
	health, err := getServiceHealth()
	if err != nil {
		log.Printf("error getting service health: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting service health"})
		return
	}

	c.JSON(http.StatusOK, health)
}

func (ct *controller) getEnrollment(c *gin.Context) {
	input := c.Param("id")
	id, err := uuid.Parse(input)
	if err != nil {
		log.Printf("error parsing path ID: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	enrollment, err := ct.service.GetEnrollment(c, id)
	if err == sql.ErrNoRows {
		log.Printf("enrollment ID %s not found", id.String())
		c.JSON(http.StatusNotFound, gin.H{"error": "enrollment not found"})
		return
	}

	if err != nil {
		log.Printf("error getting enrollment: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	filename := fmt.Sprintf(`attachment; filename=%s.csv`, enrollment.Name)
	c.Header("Content-Disposition", filename)
	c.Header("Content-Type", "text/csv")
	c.Writer.Write(enrollment.Content)
}

func (ct *controller) listEnrollment(c *gin.Context) {
	enrollments, err := ct.service.ListEnrollments(c)
	if err != nil {
		log.Printf("error getting enrollments: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, enrollments)
}

func (ct *controller) parseEnrollment(c *gin.Context) {
	var input *InputEnrollment
	if err := c.BindJSON(&input); err != nil {
		log.Printf("error parsing content body: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body content"})
		return
	}

	infoPdf, err := input.validate()
	if err != nil {
		errMsg := fmt.Sprintf("invalid body content: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	enrollment, err := ct.service.CheckEnrollmentExist(c, infoPdf.URL.String())
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error checking enrollment exist: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if enrollment != nil {
		errMsg := fmt.Sprintf("content associated with URL %s already exist. Named as %s", input.URL, enrollment.Name)
		c.JSON(http.StatusConflict, gin.H{"error": errMsg})
		return
	}

	enrollment, err = ct.service.CreateEnrollment(c, infoPdf)
	if err != nil {
		log.Printf("error creating enrollment: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   enrollment.ID,
		"name": enrollment.Name,
	})
}

func (c *controller) Start() {
	config := config.Get()

	router := gin.Default()

	if config.Server.WithCors {
		router.Use(CORSMiddleware())
	}

	router.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/ufabc-enrollment-filter"
		router.HandleContext(c)
	})
	router.Static("/ufabc-enrollment-filter", config.Server.StaticPath)

	router.GET("/health", c.health)
	router.POST("/parse-enrollment", c.parseEnrollment)
	router.GET("/enrollment/:id", c.getEnrollment)
	router.GET("/enrollment", c.listEnrollment)

	router.Run(fmt.Sprintf(":%d", config.Server.Port))
}
