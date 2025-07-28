package infraestructure

import (
	"log"
	"os"

	app_users "github.com/JosephAntony37900/Geova-back-1/Users/application"
	control_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/controllers"
	domain_users "github.com/JosephAntony37900/Geova-back-1/Users/domain/repository"
	repo_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/repository"
	routes_users "github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/routes"
	services_users"github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/services"
	"github.com/JosephAntony37900/Geova-back-1/core"

	"github.com/gin-gonic/gin"
)

// UserInfrastructure encapsula toda la infraestructura de usuarios
type UserInfrastructure struct {
	DatabaseManager *core.DatabaseManager
	UserRepo        domain_users.UserRepository
}

// NewUserInfrastructure crea e inicializa toda la infraestructura de usuarios
func NewUserInfrastructure() *UserInfrastructure {
	// Inicializar el DatabaseManager (maneja ambas conexiones)
	dbManager := core.NewDatabaseManager()
	
	// Validar que el DatabaseManager se inicializó correctamente
	if dbManager == nil {
		panic("ERROR CRÍTICO: No se pudo inicializar el DatabaseManager")
	}
	
	// Verificar estado de las conexiones
	if dbManager.LocalDB == nil || dbManager.LocalDB.DB == nil {
		panic("ERROR CRÍTICO: No se puede inicializar sin conexión a BD local")
	}
	
	// Log del estado de las conexiones
	if dbManager.RemoteDB == nil || dbManager.RemoteDB.DB == nil {
		log.Println("INFO: Iniciando en modo offline - solo BD local disponible")
		log.Println("INFO: Los datos se sincronizarán automáticamente cuando la BD remota esté disponible")
	} else {
		log.Println("INFO: Iniciando con ambas conexiones disponibles (local y remota)")
	}
	
	// Crear repositorio usando el DatabaseManager
	userRepo := repo_users.NewUserMySQLRepository(
		dbManager.LocalDB, 
		dbManager.RemoteDB,
	)
	
	return &UserInfrastructure{
		DatabaseManager: dbManager,
		UserRepo:        userRepo,
	}
}

// InitUserDependencies inicializa todas las dependencias y configura las rutas
func InitUserDependencies(engine *gin.Engine) *UserInfrastructure {
	log.Println("INFO: Inicializando infraestructura de usuarios...")
	
	// Crear infraestructura
	infrastructure := NewUserInfrastructure()
	
	// Configuración de servicios
	log.Println("INFO: Inicializando servicios de seguridad...")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("WARNING: JWT_SECRET no está configurado, usando valor por defecto")
		// En producción esto debería ser un error crítico
		// panic("JWT_SECRET no está configurado en las variables de entorno")
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

// Shutdown cierra todas las conexiones de forma limpia
func (ui *UserInfrastructure) Shutdown() {
	log.Println("INFO: Cerrando infraestructura de usuarios...")
	
	if ui.DatabaseManager != nil {
		ui.DatabaseManager.Close()
	}
	
	log.Println("INFO: Infraestructura de usuarios cerrada exitosamente")
}

// GetConnectionStatus retorna el estado de las conexiones
func (ui *UserInfrastructure) GetConnectionStatus() map[string]bool {
	status := make(map[string]bool)
	
	// Verificar conexión local
	status["local"] = false
	if ui.DatabaseManager.LocalDB != nil && ui.DatabaseManager.LocalDB.DB != nil {
		if err := ui.DatabaseManager.LocalDB.DB.Ping(); err == nil {
			status["local"] = true
		}
	}
	
	// Verificar conexión remota
	status["remote"] = false
	if ui.DatabaseManager.RemoteDB != nil && ui.DatabaseManager.RemoteDB.DB != nil {
		if err := ui.DatabaseManager.RemoteDB.DB.Ping(); err == nil {
			status["remote"] = true
		}
	}
	
	return status
}

// ReconnectRemoteDB intenta reconectar a la BD remota
func (ui *UserInfrastructure) ReconnectRemoteDB() bool {
	if ui.DatabaseManager != nil {
		ui.DatabaseManager.ReconnectRemote()
		
		// Verificar si la reconexión fue exitosa
		if ui.DatabaseManager.RemoteDB != nil && ui.DatabaseManager.RemoteDB.DB != nil {
			if err := ui.DatabaseManager.RemoteDB.DB.Ping(); err == nil {
				log.Println("INFO: Reconexión a BD remota exitosa")
				return true
			}
		}
	}
	
	log.Println("WARNING: No se pudo reconectar a la BD remota")
	return false
}

// HealthCheck verifica el estado general de la infraestructura
func (ui *UserInfrastructure) HealthCheck() map[string]interface{} {
	healthStatus := make(map[string]interface{})
	
	// Estado de conexiones
	connectionStatus := ui.GetConnectionStatus()
	healthStatus["connections"] = connectionStatus
	
	// Estado general
	healthStatus["healthy"] = connectionStatus["local"] // Mínimo requerido es la BD local
	healthStatus["mode"] = "offline"
	if connectionStatus["remote"] {
		healthStatus["mode"] = "online"
	}
	
	// Información adicional
	healthStatus["sync_enabled"] = connectionStatus["remote"]
	
	return healthStatus
}