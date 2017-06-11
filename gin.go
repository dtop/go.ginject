package ginject

import (
	"log"

	"github.com/gin-gonic/gin"
)

func DependencyInjector(i Injector) gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Set("dep", i)
		c.Next()
	}
}

func Deps(c *gin.Context) Injector {

	dep, ok := c.Get("dep")
	if !ok {
		panic("deps not present in context")
		return nil
	}

	return dep.(Injector)
}

func (i *Inj) Selfcheck(c *gin.Context) {

	for k, v := range i.deps {

		log.Println(k, v)
	}
}
