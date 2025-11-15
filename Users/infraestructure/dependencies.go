package infraestructure

import (
	"log"
	"os"

	app_users "github.com/JosephAntony37900/Geova-back-1/Users/application"
	domain_users "github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	control_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/controllers"
	repo_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/repository"
	routes_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/routes"
	services_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/services"
	"github.com/JosephAntony37900/Geova-back-1/core"

	"github.com/gin-gonic/gin"
)

type UserInfrastructure struct {
	DB       *core.Conn_MySQL
	UserRepo domain_users.UserRepository
}

func NewUserInfrastructure() *UserInfrastructure {
	// Inicializar conexión a base de datos
	db := core.NewDatabaseConnection()

	if db == nil || db.DB == nil {
		panic("ERROR CRÍTICO: No se pudo inicializar la conexión a la base de datos")
	}

	log.Println("INFO: Conexión a base de datos establecida")

	// Crear repositorio
	userRepo := repo_users.NewUserMySQLRepository(db)

	return &UserInfrastructure{
		DB:       db,
		UserRepo: userRepo,
	}
}

func InitUserDependencies(engine *gin.Engine) *UserInfrastructure {
	log.Println("INFO: Inicializando infraestructura de usuarios...")

	// Crear infraestructura
	infrastructure := NewUserInfrastructure()

	// Inicializar servicios de seguridad
	log.Println("INFO: Inicializando servicios de seguridad...")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("WARNING: JWT_SECRET no está configurado, usando valor por defecto")

	}

	bcryptService := services_users.InitBcryptService()
	jwtManager := services_users.InitTokenManager()

	if bcryptService == nil {
		panic("ERROR CRÍTICO: No se pudo inicializar el servicio de Bcrypt")
	}

	if jwtManager == nil {
		panic("ERROR CRÍTICO: No se pudo inicializar el Token Manager")
	}

	log.Println("INFO: Servicios de seguridad inicializados exitosamente")

	// Crear casos de uso
	log.Println("INFO: Inicializando casos de uso...")
	createUserUseCase := app_users.NewCreateUserUseCase(infrastructure.UserRepo, bcryptService)
	getAllUsersUseCase := app_users.NewGetUsersUseCase(infrastructure.UserRepo)
	getUserByIdUseCase := app_users.NewGetUserByIdUseCase(infrastructure.UserRepo)
	updateUserUseCase := app_users.NewUpdateUserUseCase(infrastructure.UserRepo, bcryptService)
	deleteUserUseCase := app_users.NewDeleteUserUseCase(infrastructure.UserRepo)
	loginUserUseCase := app_users.NewLoginUseCase(infrastructure.UserRepo, jwtManager, bcryptService)

	// Crear controladores
	log.Println("INFO: Inicializando controladores...")
	createUserController := control_users.NewCreateUserController(createUserUseCase)
	getAllUsersController := control_users.NewGetAllUsersController(getAllUsersUseCase)
	getUserByIdController := control_users.NewGetUserByIdUseController(getUserByIdUseCase)
	updateUserController := control_users.NewUpdateUserController(updateUserUseCase)
	deleteUserController := control_users.NewDeleteUserController(deleteUserUseCase)
	loginUserController := control_users.NewLoginUserController(loginUserUseCase)

	// Configurar rutas
	log.Println("INFO: Configurando rutas de usuarios...")
	routes_users.SetupUserRoutes(engine,
		createUserController,
		getAllUsersController,
		getUserByIdController,
		updateUserController,
		deleteUserController,
		loginUserController,
	)

	log.Println("INFO: Infraestructura de usuarios inicializada exitosamente")
	return infrastructure
}

func (ui *UserInfrastructure) Shutdown() {
	log.Println("INFO: Cerrando infraestructura de usuarios...")

	if ui.DB != nil && ui.DB.DB != nil {
		ui.DB.DB.Close()
		log.Println("INFO: Conexión a base de datos cerrada")
	}

	log.Println("INFO: Infraestructura de usuarios cerrada exitosamente")
}
