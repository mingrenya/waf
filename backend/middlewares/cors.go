// middleware/cors.go
import "github.com/gin-contrib/cors"

func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"https://yourdomain.com"},
        AllowMethods:    []string{"GET", "POST"},
        AllowHeaders:    []string{"Content-Type"},
        ExposeHeaders:   []string{"Content-Length"},
        MaxAge:          12 * time.Hour,
    })
}

