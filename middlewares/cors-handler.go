package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type CorsConfig struct {
	AllowAllOrigins  bool     `json:"allow_all_origins"`
	AllowOrigins     []string `json:"allow_origins,omitempty"`
	AllowHeaders     []string `json:"allow_headers,omitempty"`
	AllowMethods     []string `json:"allow_methods,omitempty"`
	AllowCredentials bool     `json:"allow_credentials,omitempty"`
	MaxAgeHours      int      `json:"max_age_hours,omitempty"`
	AllowWebSockets  bool     `json:"allow_web_sockets,omitempty"`
}

func CorsHandler(config *CorsConfig) gin.HandlerFunc {
	return func(gc *gin.Context) {
		ginCorsConfig := cors.DefaultConfig()
		ginCorsConfig.AllowBrowserExtensions = true
		if config.AllowAllOrigins {
			origin := gc.GetHeader("Origin")
			if "" == origin {
				origin = "*"
			}
			ginCorsConfig.AllowOrigins = []string{origin}
		} else {
			ginCorsConfig.AllowOrigins = config.AllowOrigins
		}

		if len(config.AllowHeaders) > 0 {
			for _, ah := range config.AllowHeaders {
				if "*" == ah {
					for hk, _ := range gc.Request.Header {
						if "Access-Control-Request-Headers" == hk {
							for _, ahk := range strings.Split(gc.GetHeader(hk), ",") {
								ginCorsConfig.AddAllowHeaders(ahk)
							}
							continue
						}
						ginCorsConfig.AddAllowHeaders(hk)
					}
					continue
				}
				ginCorsConfig.AddAllowHeaders(ah)
			}
		}

		if len(config.AllowMethods) > 0 {
			ginCorsConfig.AllowMethods = config.AllowMethods
		}

		ginCorsConfig.AllowCredentials = config.AllowCredentials
		if config.MaxAgeHours > 0 {
			ginCorsConfig.MaxAge = time.Duration(config.MaxAgeHours) * time.Hour
		}

		ginCorsConfig.AllowWebSockets = config.AllowWebSockets

		cors.New(ginCorsConfig)(gc)
		//gc.Next()
	}
}
