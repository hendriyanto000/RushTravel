package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/rush-travel/config"
	ep "github.com/rush-travel/endpoint"
	"github.com/rush-travel/util"
)

//class
type endpoint struct {
	Method      string
	URL         string
	Description string
}

func main() {
	config.Conf()
	router := gin.Default()
	router.Use(cors.Default())
	listEndpoint := []endpoint{
		//Paylist Endpoint
		{Method: "GET", URL: "/paylist", Description: "Get/Fetch All Paylist Data"},
		{Method: "GET", URL: "/paylist/:id", Description: "Get/Fetch Single Paylist Data by ID"},
		{Method: "POST", URL: "/paylist", Description: "Create/Insert User-Paylist Data"},
		{Method: "PUT", URL: "/paylist/:id", Description: "Edit/Update Paylist Data by ID"},
		{Method: "DELETE", URL: "/paylist/:id", Description: "Delete User-Paylist Data by ID"},
	}

	router.GET("/", func(c *gin.Context) {
		util.CallSuccessOK(c, "RushTravel-API available endpoint", listEndpoint)
	})

	router.POST("/order", ep.CreateOrder)
	router.GET("/order", ep.FetchOrder)
	router.GET("/order/:id", ep.GetOrderUser)
	router.PUT("/orderFind/:id/:orderID", ep.UpdateOrderStatus)
	//check price
	// router.GET("/orderPrice/:person", ep.OrderPrice)
	router.POST("/orderPrice/", ep.OrderPrice)

	// router.POST("/paylist", ep.Auth, ep.CreateUserPaylist)
	// router.PUT("/status/:id", ep.Auth, ep.UpdateUserPaylistStatus)
	// router.PUT("/paylist/:id", ep.Auth, ep.UpdateUserPaylist)
	// router.DELETE("/paylist/:id", ep.Auth, ep.DeleteUserPaylist)

	// router.GET("/user/:id", ep.Auth, ep.FetchSingleUser)
	router.GET("/users", ep.FetchAllUser)
	router.GET("/users/signout", ep.Logout)
	router.POST("/users/signin", ep.Login)
	router.POST("/user/signup", ep.CreateUser)
	// router.POST("/addsaldo", ep.Auth, ep.AddBalance)
	// router.PUT("/user/:id", ep.Auth, ep.UpdateUser)
	// router.PUT("/editpassword/:id", ep.Auth, ep.EditPassword)
	// router.DELETE("/user/:id", ep.Auth, ep.DeleteUser)

	// router.GET("/income", ep.Auth, ep.FetchAllIncome)
	// router.PUT("/income/:id",ep.Auth, ep.UpdateIncome)
	// router.DELETE("/income/:id", ep.Auth, ep.DeleteIncome)
	router.Run(":8000")
}
