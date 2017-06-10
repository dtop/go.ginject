# go.ginject
Go dependency injection approach using gin


## Installation

```
$ go get github.com/dtop/go.ginject
```

## Usage

1) Register to gin

```go
	
	gin  := gin.New()
	deps := ginject.New()

	gin.Use(ginglog.Logger(120))
	gin.Use(gin.Recovery())
	gin.Use(ginject.DependencyInjector(deps)) // apply 
	
	// store your services
	
	deps.Register("db", database.New())
	
	deps.RegisterLazy("redis", func() interface{} {
	    
	    return redis.New()
	})
	
```

2) Use for single objects

```go

func YourEndpoint(c *gin.Context) {

    deps := ginject.Deps(c) // alternatively: c.Get("dep").(ginject.Injector)
    
    var db database.DB
    if err := deps.Get("db", &db); err != nil {
        panic(err)
    }
    
    // use db
}

```

3) Use for structs

```go

type SomeModel struct {
    Db *database.DB `inject:"db"`
}

func YourEndpoint(c *gin.Context) {

    deps := ginject.Deps(c) // alternatively: c.Get("dep").(ginject.Injector)
    
    model := &SomeModel{}
    if err := deps.Apply(model); err != nil {
        panic(err)
    }
    
    model.Db.connect()
    defer model.Db.Close()
    // ...
}

```