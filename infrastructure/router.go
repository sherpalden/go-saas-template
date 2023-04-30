package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	gin    *gin.Engine
	RGroup *gin.RouterGroup
	env    Env
}

func NewRouter(env Env) Router {
	httpRouter := gin.New()

	//setup cors policies
	allowOrigins := []string{"*"}
	if env.Environment == "production" {
		allowOrigins = []string{"*"}
	}
	httpRouter.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
	}))

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	httpRouter.Use(gin.Recovery())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	httpRouter.Use(gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}))

	httpRouter.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"data": "go-saas-template up and running."})
	})

	return Router{
		gin:    httpRouter,
		RGroup: httpRouter.Group("/api/v1"),
		env:    env,
	}
}

func (router Router) Run() {
	if router.env.ServerPort == "" {
		router.gin.Run()
	} else {
		router.gin.Run(":" + router.env.ServerPort)
	}
}
