package main

import "github.com/gin-gonic/gin"

func main() {
	engine := gin.Default()
	engine.GET("/", handler)
	engine.GET("/login", handler)
	engine.GET("/logout", handler)
	engine.GET("/products", handler)
	engine.GET("/product/:productID", handler)
	engine.GET("/basket", handler)
	engine.GET("/about", handler)
	engine.Run(":8080")
}

func handler(c *gin.Context) {

}
