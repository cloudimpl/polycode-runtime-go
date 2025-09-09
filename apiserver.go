package runtime

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ApiServer struct {
	runtime   ClientRuntime
	ginEngine *gin.Engine
}

func (s *ApiServer) Start(port int) {
	// Create a Gin router
	s.ginEngine = gin.Default()

	s.ginEngine.GET("/v1/health", s.invokeHealthCheck)
	s.ginEngine.POST("/v1/invoke/api", s.invokeApiHandler)
	s.ginEngine.POST("/v1/invoke/service", s.invokeServiceHandler)

	// Start the Gin server
	err := s.ginEngine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to start api server: %s", err.Error())
	}
}

func (s *ApiServer) invokeHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *ApiServer) invokeApiHandler(c *gin.Context) {
	var input ApiStartEvent
	var output ApiCompleteEvent

	fmt.Println("api task received")
	if err := c.ShouldBindJSON(&input); err != nil {
		output = ErrorToApiComplete(ErrInternal.Wrap(err))
		fmt.Printf("api task failed %s\n", err.Error())
	} else {
		output = s.runtime.RunApi(c, input)
		fmt.Println("api task success")
	}

	c.JSON(http.StatusOK, output)
}

func (s *ApiServer) invokeServiceHandler(c *gin.Context) {
	var input ServiceStartEvent
	var output ServiceCompleteEvent

	fmt.Println("service task received")
	if err := c.ShouldBindJSON(&input); err != nil {
		output = ErrorToServiceComplete(ErrInternal.Wrap(err), "")
		fmt.Printf("service task failed %s\n", err.Error())
	} else {
		output = s.runtime.RunService(c, input)
		fmt.Println("service task success")
	}

	c.JSON(http.StatusOK, output)
}
