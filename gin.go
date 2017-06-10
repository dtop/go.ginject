package ginject

import "github.com/gin-gonic/gin"

func DependencyInjector(i Injector) gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Set("dep", i)
		c.Next()
	}
}

func Deps(c *gin.Context) Injector {

	dep, ok := c.Get("dep")
	if !ok { return nil }

	return dep.(Injector)
}
