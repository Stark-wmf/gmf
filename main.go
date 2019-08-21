package main

import (
	"gmf/gmf"
)
// LoginJSON .
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

func getInfo(c *gmf.Context) string {
	name := c.Get("innerName")
	message := "welcome " + name.(string) // 前提是知道这个是string类型
	return message
}

func Submitfunc(c *gmf.Context) {
	c.String(200, "submit")

}

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




	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}