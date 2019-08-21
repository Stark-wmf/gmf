# 快速启动-QUICK START

```
func main() {
   r := gmf.New()
   r.Use()
   r.GET("/hello")
 // Listen and server on 0.0.0.0:8080
   r.Run(":8080")
}
```

## 路由方法使用-HOW TO USE METHOD

```
func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
```

## 路径中的参数-Parameters in path

```
func main() {
   r := gmf.New()
   r.Use(gmf.Logger(),gmf.Recovery())
   // Simple group: v1
   g := r.Group("/1234")
   {
      g.GET("/login/:name", Loginfunc)
      g.GET("/login/:bool", Loginfunc)
      g.POST("/passwordlogin", PasswordLoginfunc)
      g.GET("/submit", Submitfunc)
   }

    r.GET("/login")
   r.POST("/post")


   // Listen and server on 0.0.0.0:8080
   r.Run(":8080")
}
```

## 路由组-Grouping routes

```
func main() {
   r := gmf.New()
   r.Use(gmf.Logger(),gmf.Recovery())
   // Simple group: v1
   g := r.Group("/1234")
   {
      g.GET("/login/:name", Loginfunc)
      g.GET("/submit", Submitfunc)
   }
   r.POST("/post")

   // Listen and server on 0.0.0.0:8080
   r.Run(":8080")
}
```

## 使用中间件-Using middleware

```
func main() {
   r := gmf.New()
   //使用Use来添加中间件
   r.Use(gmf.Logger())
   r.Use(gin.Recovery())
   // Simple group: v1
   g := r.Group("/1234")
   {
      g.GET("/login/:name", Loginfunc)
      g.GET("/login/:bool", Loginfunc)
      g.POST("/passwordlogin", PasswordLoginfunc)
      g.GET("/submit", Submitfunc)
   }

    r.GET("/login")
   r.POST("/post")


   // Listen and server on 0.0.0.0:8080
   r.Run(":8080")
}
```

## 回复JSON-RESPONSE JSON

```
type LoginJSON struct {
   User     string `json:"user" `
   Password string `json:"password" `
}

func PasswordLoginfunc(c *gmf.Context) {
   var json LoginJSON

      if json.User == "wmf" && json.Password == "191513" {
            c.JSON(200, gmf.H{"status": "you are logged in"})
      } else {
         c.JSON(401, gmf.H{"status": "unauthorized"})
      }
   }
```

## 回复STRING-RESOPNSE STRING

```
func Loginfunc(c *gmf.Context) {
   name := c.Params.ByName("name")
   c.Set("innerName", name)
   message := getInfo(c)
   //Alarm的示范用法
   if message==""{
      gmf.Alarm("INFO","mes不存在")
   }
   if message=="sss"{
      gmf.Alarm("WARN","")
   }
   c.String(200, message)

}
```