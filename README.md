# Geova Backend

API REST desarrollada en Go para la gestión de usuarios y proyectos, implementando Arquitectura Hexagonal (Ports and Adapters) con Clean Architecture.

## Tabla de Contenidos

- [Arquitectura](#arquitectura)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Módulos](#módulos)
- [Tecnologías](#tecnologías)
- [Configuración](#configuración)
- [Instalación](#instalación)
- [Ejecución](#ejecución)
- [API Endpoints](#api-endpoints)
- [Base de Datos](#base-de-datos)

## Arquitectura

El proyecto implementa **Arquitectura Hexagonal** combinada con principios de **Clean Architecture**, organizando el código en capas claramente definidas que separan las responsabilidades y facilitan el mantenimiento y testing.

### Principios Aplicados

- **Separación de Responsabilidades**: Cada capa tiene una responsabilidad única y bien definida.
- **Inversión de Dependencias**: Las capas internas no conocen las implementaciones de las capas externas.
- **Independencia de Frameworks**: La lógica de negocio no depende de frameworks específicos.
- **Testabilidad**: Las interfaces permiten fácil implementación de mocks para testing.

### Capas de la Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                      │
│  (Controllers, Routes, Repositories, Adapters, Services)    │
│  - HTTP Handlers                                             │
│  - Database Implementation                                   │
│  - External Services (Cloudinary, JWT, Bcrypt)              │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   Application Layer                          │
│                    (Use Cases)                               │
│  - Business Logic Orchestration                              │
│  - Application-specific operations                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                     Domain Layer                             │
│            (Entities, Interfaces)                            │
│  - Core Business Entities                                    │
│  - Repository Interfaces                                     │
│  - Service Interfaces                                        │
└─────────────────────────────────────────────────────────────┘
```

#### 1. Domain Layer (Capa de Dominio)

Contiene la lógica de negocio pura y las entidades del sistema. Esta capa no tiene dependencias externas.

**Componentes:**
- **Entities**: Estructuras de datos que representan los conceptos del negocio (`User`, `Project`)
- **Repository Interfaces**: Contratos que definen operaciones de persistencia
- **Service Interfaces**: Contratos para servicios externos (autenticación, encriptación, almacenamiento)

**Ejemplo:**
```go
type User struct {
    Id        int
    Username  string
    Nombre    string
    Apellidos string
    Email     string
    Password  string
}

type UserRepository interface {
    Save(user entities.User) error
    FindById(id int) (*entities.User, error)
    FindAll() ([]entities.User, error)
    Update(user entities.User) error
    Delete(id int) error
}
```

#### 2. Application Layer (Capa de Aplicación)

Contiene los casos de uso del sistema. Orquesta el flujo de datos entre la capa de dominio y la infraestructura.

**Responsabilidades:**
- Implementar casos de uso específicos de la aplicación
- Coordinar operaciones entre repositorios y servicios
- Validar reglas de negocio
- Gestionar transacciones

**Ejemplos de Casos de Uso:**
- `CreateUserUseCase`: Crea un usuario, encripta su contraseña
- `LoginUseCase`: Autentica usuarios y genera tokens JWT
- `CreateProjectUseCase`: Crea proyectos y maneja la subida de imágenes a Cloudinary
- `UpdateProjectUseCase`: Actualiza proyectos, gestiona cambio de imágenes

#### 3. Infrastructure Layer (Capa de Infraestructura)

Contiene todas las implementaciones concretas y detalles técnicos. Se comunica con el mundo exterior.

**Componentes:**

- **Controllers**: Manejan las peticiones HTTP, validan entrada, invocan casos de uso
- **Routes**: Configuran los endpoints HTTP y middlewares
- **Repositories**: Implementan las interfaces del dominio para persistencia en MySQL
- **Adapters**: Adaptan servicios externos a las interfaces del dominio
- **Services**: Implementaciones de servicios (JWT, Bcrypt, Cloudinary)
- **Dependencies**: Inyección de dependencias y configuración de módulos

### Core Layer

Capa transversal que proporciona utilidades y configuraciones comunes:

- **Database Configuration**: Gestión de conexiones a MySQL
- **CORS Setup**: Configuración de políticas de origen cruzado
- **Connection Pool**: Pool de conexiones optimizado

## Estructura del Proyecto

```
Geova-back-1/
│
├── main.go                      # Punto de entrada de la aplicación
├── go.mod                       # Dependencias del proyecto
├── go.sum                       # Checksums de dependencias
│
├── core/                        # Configuraciones y utilidades comunes
│   ├── cors.go                  # Configuración CORS
│   ├── dtabase_config.go        # Configuración de base de datos
│   └── mysql.go                 # Pool de conexiones MySQL
│
├── Users/                       # Módulo de Usuarios
│   ├── domain/
│   │   ├── entities/
│   │   │   └── users.go         # Entidad User
│   │   ├── repository/
│   │   │   └── users_repository.go    # Interface UserRepository
│   │   └── services/
│   │       ├── bcrypt.service.go      # Interface de encriptación
│   │       └── token_manager.go       # Interface de gestión JWT
│   │
│   ├── application/             # Casos de uso
│   │   ├── createUsers_useCase.go
│   │   ├── loginUsers_useCase.go
│   │   ├── getUserById_useCase.go
│   │   ├── getUsers_useCase.go
│   │   ├── updateUsers_useCase.go
│   │   └── deleteUsers_useCase.go
│   │
│   └── infraestructure/
│       ├── dependencies.go      # Inyección de dependencias
│       ├── controllers/         # Controladores HTTP
│       │   ├── createUsers_Controllers.go
│       │   ├── loginController.go
│       │   ├── getAllUsers_Contrroller.go
│       │   ├── getUserById_Controller.go
│       │   ├── updateUsers_Contoller.go
│       │   └── deleteUser_Controller.go
│       ├── repository/
│       │   └── users_repo_mysql.go    # Implementación MySQL
│       ├── routes/
│       │   └── users_routers.go       # Rutas HTTP
│       ├── adapters/
│       │   ├── bcript_adapter.service.go
│       │   └── jwt_manager.go
│       └── services/
│           └── Middlewares_Auth.go    # Middleware de autenticación
│
└── Projects/                    # Módulo de Proyectos
    ├── domain/
    │   ├── entities/
    │   │   └── projects.go      # Entidad Project
    │   ├── repository/
    │   │   └── repository_projects.go  # Interface ProjectRepository
    │   └── services/
    │       └── cloudinary_service.go   # Interface de almacenamiento
    │
    ├── application/             # Casos de uso
    │   ├── createProject_usecase.go
    │   ├── getAllProjects_useCase.go
    │   ├── getByIdProject_usecase.go
    │   ├── getProjectByName_useCase.go
    │   ├── getProjectByCategory_useCase.go
    │   ├── getProjectByDate_useCase.go
    │   ├── getProjectByIdUser.go
    │   ├── updateProjects_usecase.go
    │   └── deleteProject_useCase.go
    │
    └── infraestructure/
        ├── projects_dependencies.go
        ├── controllers/         # Controladores HTTP
        │   ├── createProjects_controller.go
        │   ├── getAllProjects_controller.go
        │   ├── getByIdProjects_controller.go
        │   ├── getByNameProject_controller.go
        │   ├── getProjectByCategory_controller.go
        │   ├── getProjectsByDate_controller.go
        │   ├── getProjectsByIdUser_controller.go
        │   ├── updateProjectById_controller.go
        │   └── deleteProjects_controller.go
        ├── repository/
        │   └── projects_repo_mysql.go   # Implementación MySQL
        ├── routes/
        │   └── proyects_routes.go       # Rutas HTTP
        └── services/
            └── adapters/
                └── cloudinary_adapter.go # Adaptador Cloudinary
```

## Módulos

### Módulo Users

Gestiona la autenticación, autorización y operaciones CRUD de usuarios.

**Características:**
- Registro de usuarios con encriptación de contraseñas (Bcrypt)
- Autenticación mediante JWT (JSON Web Tokens)
- Gestión completa de perfiles de usuario
- Middleware de autenticación para rutas protegidas

**Entidad User:**
```go
type User struct {
    Id        int
    Username  string
    Nombre    string
    Apellidos string
    Email     string
    Password  string  // Hash Bcrypt
}
```

**Casos de Uso:**
1. **CreateUser**: Registra un nuevo usuario
   - Valida que el email no esté registrado
   - Encripta la contraseña con Bcrypt
   - Guarda el usuario en la base de datos

2. **Login**: Autentica un usuario
   - Busca usuario por email
   - Verifica contraseña con Bcrypt
   - Genera token JWT válido por 24 horas

3. **GetUsers**: Obtiene lista de todos los usuarios

4. **GetUserById**: Obtiene un usuario específico por ID

5. **UpdateUser**: Actualiza información de usuario
   - Verifica existencia del usuario
   - Re-encripta contraseña si cambió
   - Valida unicidad de email

6. **DeleteUser**: Elimina un usuario del sistema

### Módulo Projects

Gestiona proyectos con geolocalización e imágenes almacenadas en Cloudinary.

**Características:**
- CRUD completo de proyectos
- Integración con Cloudinary para almacenamiento de imágenes
- Búsquedas por nombre, categoría, fecha y usuario
- Soporte de geolocalización (latitud/longitud)

**Entidad Project:**
```go
type Project struct {
    Id             int
    NombreProyecto string
    Fecha          string
    Categoria      string
    Descripcion    string
    Img            string  // URL de Cloudinary
    Lat            float64
    Lng            float64
    UserId         int
}
```

**Casos de Uso:**
1. **CreateProject**: Crea un nuevo proyecto
   - Sube imagen a Cloudinary
   - Almacena URL de la imagen
   - Guarda proyecto en base de datos

2. **GetAllProjects**: Obtiene lista de proyectos ordenados por ID descendente

3. **GetProjectById**: Obtiene un proyecto específico

4. **GetProjectsByName**: Busca proyectos por nombre (LIKE)

5. **GetProjectsByCategory**: Filtra proyectos por categoría

6. **GetProjectsByDate**: Filtra proyectos por fecha

7. **GetProjectsByUserId**: Obtiene proyectos de un usuario específico

8. **UpdateProject**: Actualiza un proyecto
   - Actualiza imagen en Cloudinary si cambió
   - Modifica datos del proyecto

9. **DeleteProject**: Elimina un proyecto

## Tecnologías

### Framework y Librerías

- **Go 1.24.4**: Lenguaje de programación
- **Gin 1.10.1**: Framework web HTTP
- **gin-contrib/cors**: Middleware CORS
- **go-sql-driver/mysql 1.9.3**: Driver MySQL

### Seguridad

- **golang.org/x/crypto**: Encriptación Bcrypt
- **golang-jwt/jwt/v4**: Generación y validación de JWT
- **dgrijalva/jwt-go**: Soporte JWT adicional

### Base de Datos

- **MySQL**: Base de datos relacional

### Servicios Externos

- **Cloudinary v2.11.0**: Almacenamiento y gestión de imágenes

### Utilidades

- **godotenv 1.5.1**: Gestión de variables de entorno

## Configuración

### Variables de Entorno

Crear un archivo `.env` en la raíz del proyecto:

```env
# Base de Datos Remota
REMOTE_DB_HOST=your-database-host.com
REMOTE_DB_USER=your-username
REMOTE_DB_PASS=your-password
REMOTE_DB_SCHEMA=your-database-name
REMOTE_DB_PORT=3306

# JWT
JWT_SECRET=your-jwt-secret-key-here

# Cloudinary
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret

# CORS (opcional)
ALLOWED_ORIGIN=https://your-frontend-domain.com

# Rate Limiting para proyectos
PROJECTS_WRITE_RATE_LIMIT=you-valor-of-configuration-here
PROJECTS_WRITE_BURST_LIMIT=you-valor-of-configuration-here
PROJECTS_READ_RATE_LIMIT=you-valor-of-configuration-here
PROJECTS_READ_BURST_LIMIT=you-valor-of-configuration-here
PROJECTS_QUERY_RATE_LIMIT=you-valor-of-configuration-here
PROJECTS_QUERY_BURST_LIMIT=you-valor-of-configuration-here
PROJECTS_RATE_LIMIT_TTL=you-valor-of-configuration-here
PROJECTS_RATE_LIMIT_CLEANUP=5m

# Rate Limiting para usuarios
USERS_LOGIN_RATE_LIMIT=you-valor-of-configuration-here
USERS_LOGIN_BURST_LIMIT=you-valor-of-configuration-here
USERS_REGISTER_RATE_LIMIT=you-valor-of-configuration-here
USERS_REGISTER_BURST_LIMIT=you-valor-of-configuration-here
USERS_MODIFY_RATE_LIMIT=you-valor-of-configuration-here
USERS_MODIFY_BURST_LIMIT=you-valor-of-configuration-here
USERS_READ_RATE_LIMIT=you-valor-of-configuration-here
USERS_READ_BURST_LIMIT=you-valor-of-configuration-here
USERS_RATE_LIMIT_TTL=you-valor-of-configuration-here
USERS_RATE_LIMIT_CLEANUP=you-valor-of-configuration-here
```

### CORS

La configuración CORS permite peticiones desde los siguientes orígenes:

- `https://www.geova.pro`
- `https://geova.pro`
- `http://localhost:3000` (desarrollo)
- `http://localhost:5173` (desarrollo Vite)
- `http://localhost:4200` (desarrollo Angular)
- Origen adicional vía variable `ALLOWED_ORIGIN`

**Headers permitidos:**
- Origin
- Content-Type
- Authorization
- Accept
- X-Requested-With

**Métodos HTTP permitidos:**
- GET
- POST
- PUT
- PATCH
- DELETE
- OPTIONS

## Instalación

### Prerrequisitos

- Go 1.24 o superior
- MySQL 8.0 o superior
- Cuenta en Cloudinary (para almacenamiento de imágenes)

### Pasos

1. Clonar el repositorio:
```bash
git clone https://github.com/JosephAntony37900/Geova-back-1.git
cd Geova-back-1
```

2. Instalar dependencias:
```bash
go mod download
```

3. Configurar variables de entorno:
```bash
cp .env.example .env
# Editar .env con tus credenciales
```

4. Crear base de datos:
```sql
CREATE DATABASE geova_db;
USE geova_db;

-- Tabla users
CREATE TABLE users (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    Username VARCHAR(100) NOT NULL,
    Nombre VARCHAR(100) NOT NULL,
    Apellidos VARCHAR(100) NOT NULL,
    Email VARCHAR(150) NOT NULL UNIQUE,
    Password VARCHAR(255) NOT NULL,
    INDEX idx_email (Email)
);

-- Tabla projects
CREATE TABLE projects (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    NombreProyecto VARCHAR(200) NOT NULL,
    Fecha VARCHAR(50) NOT NULL,
    Categoria VARCHAR(100) NOT NULL,
    Descripcion TEXT,
    Img VARCHAR(500),
    Lat DECIMAL(10, 8),
    Lng DECIMAL(11, 8),
    user_id INT NOT NULL,
    INDEX idx_categoria (Categoria),
    INDEX idx_fecha (Fecha),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(Id) ON DELETE CASCADE
);
```

## Ejecución

### Desarrollo

```bash
go run main.go
```

### Producción

```bash
# Compilar
go build -o geova-backend

# Ejecutar
./geova-backend
```

El servidor se iniciará en `0.0.0.0:8000`

## API Endpoints

### Usuarios

#### Registro
```http
POST /users
Content-Type: application/json

{
    "username": "johndoe",
    "nombre": "John",
    "apellidos": "Doe",
    "email": "john@example.com",
    "password": "securepassword123"
}
```

#### Login
```http
POST /users/login
Content-Type: application/json

{
    "email": "john@example.com",
    "password": "securepassword123"
}

Response:
{
    "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Obtener Usuarios (Protegido)
```http
GET /users
Authorization: Bearer {token}
```

#### Obtener Usuario por ID (Protegido)
```http
GET /users/{id}
Authorization: Bearer {token}
```

#### Actualizar Usuario (Protegido)
```http
PUT /users/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
    "username": "johndoe_updated",
    "nombre": "John",
    "apellidos": "Doe Updated",
    "email": "john.updated@example.com",
    "password": "newsecurepassword123"
}
```

#### Eliminar Usuario (Protegido)
```http
DELETE /users/{id}
Authorization: Bearer {token}
```

### Proyectos

#### Crear Proyecto (Protegido)
```http
POST /projects
Authorization: Bearer {token}
Content-Type: multipart/form-data

nombreProyecto: Proyecto Ejemplo
fecha: 2025-11-15
categoria: Tecnología
descripcion: Descripción del proyecto
img: [archivo de imagen]
lat: 19.432608
lng: -99.133209
userId: 1
```

#### Obtener Todos los Proyectos
```http
GET /projects
```

#### Obtener Proyecto por ID
```http
GET /projects/{id}
```

#### Buscar Proyectos por Nombre
```http
GET /projects/search?name={nombre}
```

#### Buscar Proyectos por Categoría
```http
GET /projects/category/{categoria}
```

#### Buscar Proyectos por Fecha
```http
GET /projects/date/{fecha}
```

#### Obtener Proyectos por Usuario
```http
GET /projects/user/{userId}
```

#### Actualizar Proyecto (Protegido)
```http
PUT /projects/{id}
Authorization: Bearer {token}
Content-Type: multipart/form-data

nombreProyecto: Proyecto Actualizado
fecha: 2025-11-16
categoria: Tecnología
descripcion: Nueva descripción
img: [nuevo archivo de imagen opcional]
lat: 19.432608
lng: -99.133209
userId: 1
```

#### Eliminar Proyecto (Protegido)
```http
DELETE /projects/{id}
Authorization: Bearer {token}
```

## Base de Datos

### Esquema de Base de Datos

#### Tabla: users
```sql
CREATE TABLE users (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    Username VARCHAR(100) NOT NULL,
    Nombre VARCHAR(100) NOT NULL,
    Apellidos VARCHAR(100) NOT NULL,
    Email VARCHAR(150) NOT NULL UNIQUE,
    Password VARCHAR(255) NOT NULL,
    INDEX idx_email (Email)
);
```

**Campos:**
- `Id`: Identificador único autoincremental
- `Username`: Nombre de usuario
- `Nombre`: Nombre(s) del usuario
- `Apellidos`: Apellido(s) del usuario
- `Email`: Correo electrónico único
- `Password`: Contraseña hasheada con Bcrypt

#### Tabla: projects
```sql
CREATE TABLE projects (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    NombreProyecto VARCHAR(200) NOT NULL,
    Fecha VARCHAR(50) NOT NULL,
    Categoria VARCHAR(100) NOT NULL,
    Descripcion TEXT,
    Img VARCHAR(500),
    Lat DECIMAL(10, 8),
    Lng DECIMAL(11, 8),
    user_id INT NOT NULL,
    INDEX idx_categoria (Categoria),
    INDEX idx_fecha (Fecha),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(Id) ON DELETE CASCADE
);
```

**Campos:**
- `Id`: Identificador único autoincremental
- `NombreProyecto`: Nombre del proyecto
- `Fecha`: Fecha del proyecto (formato string)
- `Categoria`: Categoría del proyecto
- `Descripcion`: Descripción detallada
- `Img`: URL de la imagen en Cloudinary
- `Lat`: Latitud (coordenada geográfica)
- `Lng`: Longitud (coordenada geográfica)
- `user_id`: ID del usuario creador (clave foránea)

**Relaciones:**
- Un usuario puede tener múltiples proyectos (1:N)
- La eliminación de un usuario elimina sus proyectos (CASCADE)

### Índices

Los índices están optimizados para las consultas más frecuentes:
- `idx_email` en users para búsquedas y login
- `idx_categoria` en projects para filtrado por categoría
- `idx_fecha` en projects para filtrado por fecha
- `idx_user_id` en projects para consultas de proyectos por usuario

## Flujo de Datos

### Ejemplo: Creación de Usuario

```
Cliente HTTP
    │
    ├──> POST /users (JSON)
    │
    ▼
[CreateUserController]
    │
    ├──> Parsea request
    ├──> Valida datos
    │
    ▼
[CreateUserUseCase]
    │
    ├──> Verifica email único (UserRepository.FindByEmail)
    ├──> Encripta password (BcryptService.HashPassword)
    │
    ▼
[UserRepository]
    │
    ├──> INSERT INTO users...
    │
    ▼
[MySQL Database]
    │
    └──> Usuario guardado
```

### Ejemplo: Login y Autenticación

```
Cliente HTTP
    │
    ├──> POST /users/login (email, password)
    │
    ▼
[LoginController]
    │
    ▼
[LoginUseCase]
    │
    ├──> Busca usuario (UserRepository.FindByEmail)
    ├──> Verifica password (BcryptService.ComparePassword)
    ├──> Genera JWT (TokenManager.GenerateToken)
    │
    ▼
Cliente recibe Token
    │
    ├──> Guarda token
    │
    └──> Usa token en peticiones protegidas
         │
         ▼
    [AuthMiddleware]
         │
         ├──> Valida token JWT
         ├──> Extrae claims
         │
         ▼
    [Controlador protegido]
```

### Ejemplo: Creación de Proyecto con Imagen

```
Cliente HTTP
    │
    ├──> POST /projects (multipart/form-data)
    │
    ▼
[AuthMiddleware]
    │
    ├──> Valida token JWT
    │
    ▼
[CreateProjectController]
    │
    ├──> Parsea multipart form
    ├──> Extrae archivo de imagen
    │
    ▼
[CreateProjectUseCase]
    │
    ├──> Sube imagen (CloudinaryService.UploadImage)
    │    │
    │    └──> [Cloudinary API]
    │         └──> Retorna URL de imagen
    │
    ├──> Crea entidad Project con URL
    │
    ▼
[ProjectRepository]
    │
    ├──> INSERT INTO projects...
    │
    ▼
[MySQL Database]
    │
    └──> Proyecto guardado con URL de imagen
```

## Patrones de Diseño Utilizados

### 1. Dependency Injection

Todas las dependencias se inyectan a través de constructores, facilitando testing y mantenibilidad.

```go
func NewCreateUserUseCase(
    repo repository.UserRepository,
    bcrypt services.BcryptService,
) *CreateUserUseCase {
    return &CreateUserUseCase{
        repo:   repo,
        bcrypt: bcrypt,
    }
}
```

### 2. Repository Pattern

Abstrae el acceso a datos mediante interfaces, permitiendo cambiar la implementación sin afectar la lógica de negocio.

```go
type UserRepository interface {
    Save(user entities.User) error
    FindById(id int) (*entities.User, error)
    // ...
}
```

### 3. Adapter Pattern

Adapta servicios externos a las interfaces del dominio.

```go
type CloudinaryAdapter struct {
    client *cloudinary.Cloudinary
}

func (c *CloudinaryAdapter) UploadImage(file multipart.File) (string, error) {
    // Implementación específica de Cloudinary
}
```

### 4. Use Case Pattern

Encapsula la lógica de aplicación en casos de uso independientes y reutilizables.

### 5. Middleware Pattern

Intercepta peticiones HTTP para cross-cutting concerns como autenticación.

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Validar token
        // Abortar o continuar
    }
}
```

## Seguridad

### Autenticación

- **JWT (JSON Web Tokens)**: Tokens firmados con HMAC-SHA256
- **Duración**: 24 horas
- **Claims**: Incluyen ID y email del usuario

### Encriptación

- **Bcrypt**: Algoritmo de hashing para contraseñas
- **Cost Factor**: 10 (configurable)
- **Salt**: Generado automáticamente por Bcrypt

### Protección de Rutas

Middleware de autenticación valida JWT en cada petición a rutas protegidas.

### CORS

Configuración restrictiva que solo permite orígenes específicos en producción.

## Mejores Prácticas Implementadas

1. **Separación de Responsabilidades**: Cada componente tiene una única responsabilidad
2. **Código Limpio**: Nombres descriptivos, funciones pequeñas
3. **Error Handling**: Manejo consistente de errores en todas las capas
4. **Validación**: Validación de datos en controladores y casos de uso
5. **Logging**: Logs informativos para debugging y monitoreo
6. **Connection Pooling**: Pool de conexiones MySQL optimizado
7. **Prepared Statements**: Prevención de SQL Injection
8. **HTTPS Ready**: Preparado para producción con HTTPS

## Licencia

Este proyecto es privado y confidencial.

## Autor

Joseph Antony - JosephAntony37900
