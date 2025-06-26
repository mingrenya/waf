// middleware/ratelimit.go
import "github.com/ulule/limiter/v3"

func RateLimitMiddleware() gin.HandlerFunc {
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100,
    }
    store := memory.NewStore()
    limiter := limiter.New(store, rate)
    
    return func(c *gin.Context) {
        context, err := limiter.Get(c, "global")
        if err != nil {
            c.AbortWithStatus(500)
            return
        }
        
        if context.Reached {
            c.AbortWithStatusJSON(429, gin.H{
                "error": "Too many requests",
            })
            return
        }
        c.Next()
    }
}

