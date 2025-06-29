package infraestructure

import (
	_"database/sql"
	"os"

	app_users "github.com/JosephAntony37900/Geova-back-1/Users/application"
	control_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/controllers"
	repo_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/repository"
	routes_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/routes"
	"github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/services"
	"github.com/JosephAntony37900/Geova-back-1/core"
	"github.com/gin-gonic/gin"
)

func InitUserDependencies(engine *gin.Engine, conn *core.Conn_MySQL) {
	// Configuración de servicios
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET no está configurado en las variables de entorno")
	}

	bcryptService := service.InitBcryptService()
    jwtManager := service.InitTokenManager()

	userRepo := repo_users.NewUserMySQLRepository(conn)
	//3
	createUserUseCase := app_users.NewCreateUserUseCase(userRepo, bcryptService)
	getAllUsersUseCase := app_users.NewGetUsersUseCase(userRepo)
	getUserByIdUseCase := app_users.NewGetUserByIdUseCase(userRepo)
	upateUserUseCase := app_users.NewUpdateUserUseCase(userRepo, bcryptService)
	deleteUserUseCase := app_users.NewDeleteUserUseCase(userRepo)
	loginUserUsecas := app_users.NewLoginUseCase(userRepo, jwtManager, bcryptService )

	createUserController := control_users.NewCreateUserController(createUserUseCase)
	getAllUsersController := control_users.NewGetAllUsersController(getAllUsersUseCase)
	getUserByIdUserController := control_users.NewGetUserByIdUseController(getUserByIdUseCase)
	updateUserController := control_users.NewUpdateUserController(upateUserUseCase)
	deleteUserController := control_users.NewDeleteUserController(deleteUserUseCase)
	loginUserController := control_users.NewLoginUserController(loginUserUsecas)
	
	routes_users.SetupUserRoutes(engine, createUserController, getAllUsersController, getUserByIdUserController, updateUserController, deleteUserController, loginUserController)

}