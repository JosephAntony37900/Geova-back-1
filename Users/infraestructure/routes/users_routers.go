package routes

import (
	"github.com/gin-gonic/gin"
	_"os"
	"github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/controllers"
	_"github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/services"
)

func SetupUserRoutes(r *gin.Engine, createUserController *controllers.CreateUserController,
		getUsersController *controllers.GetAllUsersController,
		getUsersControllerById *controllers.GetUserByIdController,
		updateUserController *controllers.UpdateUserController,
		deleteUserController *controllers.DeleteUserController,
		loginUserController *controllers.LoginUserController,
		syncronizeUsersController *controllers.SyncUsersController) {
	//jwtSecret := os.Getenv("JWT_SECRET")

	r.POST("/users", createUserController.Execute)
	r.GET("/users", getUsersController.Execute)
	r.GET("/users/:id", getUsersControllerById.Execute)
	r.PUT("/users/:id", updateUserController.Execute)
	r.DELETE("/users/:id", deleteUserController.Execute)
	r.POST("/users/login", loginUserController.Execute )
	r.POST("/sync/users", syncronizeUsersController.Execute)
}