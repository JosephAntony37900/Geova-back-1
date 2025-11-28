// geova-back-1/Users/application/users_usecase_test.go
package application

import (
	"errors"
	"testing"

	"github.com/JosephAntony37900/Geova-back-1/Users/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Users/infraestructure/adapters"
)

// ============================================================================
// MOCKS para testing de Users
// ============================================================================

// MockUserRepository simula el repositorio de usuarios
// Implementa TODOS los métodos de la interfaz UserRepository
func TestDummy(t *testing.T) {
	// Esta función vacía permite que go test reconozca los benchmarks
}

type MockUserRepository struct {
	users map[string]*entities.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*entities.User),
	}
}

func (m *MockUserRepository) Save(user entities.User) error {
	m.users[user.Email] = &user
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*entities.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, errors.New("usuario no encontrado")
}

func (m *MockUserRepository) FindAll() ([]entities.User, error) {
	users := make([]entities.User, 0)
	for _, u := range m.users {
		users = append(users, *u)
	}
	return users, nil
}

func (m *MockUserRepository) FindById(id int) (*entities.User, error) {
	for _, u := range m.users {
		if u.Id == id {
			return u, nil
		}
	}
	return nil, errors.New("usuario no encontrado")
}

func (m *MockUserRepository) Update(user entities.User) error {
	m.users[user.Email] = &user
	return nil
}

func (m *MockUserRepository) Delete(id int) error {
	for email, u := range m.users {
		if u.Id == id {
			delete(m.users, email)
			return nil
		}
	}
	return errors.New("usuario no encontrado")
}

// MockTokenManager simula el generador de tokens JWT
// Implementa la interfaz services.TokenManager
type MockTokenManager struct{}

func (m *MockTokenManager) GenerateToken(userId int) (string, error) {
	return "mock.jwt.token.xyz", nil
}

func (m *MockTokenManager) ValidateToken(token string) (bool, map[string]interface{}, error) {
	// Simula validación exitosa
	if token == "" {
		return false, nil, errors.New("token vacío")
	}
	
	// Retorna claims simulados
	claims := map[string]interface{}{
		"userId": 1,
		"exp":    1234567890,
	}
	
	return true, claims, nil
}

// ============================================================================
// BENCHMARKS - CreateUser (ANTES Y DESPUÉS de optimizaciones)
// ============================================================================

// BenchmarkCreateUser mide el rendimiento COMPLETO de la creación de usuario
// EJECUTAR ANTES: en rama main sin optimizaciones
// EJECUTAR DESPUÉS: en rama PR con optimizaciones de concurrencia
func BenchmarkCreateUser(b *testing.B) {
	mockRepo := NewMockUserRepository()
	bcryptService := adapters.NewBcrypt() // cost=12
	useCase := NewCreateUserUseCase(mockRepo, bcryptService)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		testUser := entities.User{
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "Test123!@#",
			Nombre:    "Test",
			Apellidos: "User",
		}
		delete(mockRepo.users, testUser.Email)
		b.StartTimer()

		_, err := useCase.Execute(testUser)
		if err != nil {
			b.Fatalf("Error en CreateUser: %v", err)
		}
	}
}

// BenchmarkCreateUser_HashingOnly mide EXCLUSIVAMENTE el hashing de Bcrypt
// Este es el cuello de botella principal (50-100ms por operación)
// Si optimizas con goroutines, este número NO debería cambiar
func BenchmarkCreateUser_HashingOnly(b *testing.B) {
	bcryptService := adapters.NewBcrypt()
	password := "Test123!@#"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := bcryptService.HashPassword(password)
		if err != nil {
			b.Fatalf("Error en HashPassword: %v", err)
		}
	}
}


// ============================================================================
// BENCHMARKS - Login (ANTES Y DESPUÉS de optimizaciones)
// ============================================================================

// BenchmarkLogin mide el rendimiento COMPLETO del login
// EJECUTAR ANTES: en rama main sin optimizaciones
// EJECUTAR DESPUÉS: en rama PR con optimizaciones de concurrencia
func BenchmarkLogin(b *testing.B) {
	mockRepo := NewMockUserRepository()
	bcryptService := adapters.NewBcrypt()
	mockJWT := &MockTokenManager{}
	useCase := NewLoginUseCase(mockRepo, mockJWT, bcryptService)

	// Pre-crear un usuario con contraseña hasheada
	hashedPassword, _ := bcryptService.HashPassword("Test123!@#")
	mockRepo.users["test@example.com"] = &entities.User{
		Id:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		Nombre:    "Test",
		Apellidos: "User",
	}

	input := LoginInput{
		Email:    "test@example.com",
		Password: "Test123!@#",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := useCase.Execute(input)
		if err != nil {
			b.Fatalf("Error en Login: %v", err)
		}
	}
}

// BenchmarkLogin_CompareOnly mide EXCLUSIVAMENTE la comparación de Bcrypt
// Este es el cuello de botella del login (50-100ms por operación)
// Si optimizas con goroutines, este número NO debería cambiar
func BenchmarkLogin_CompareOnly(b *testing.B) {
	bcryptService := adapters.NewBcrypt()
	password := "Test123!@#"

	hashedPassword, err := bcryptService.HashPassword(password)
	if err != nil {
		b.Fatalf("Error generando hash: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bcryptService.ComparePasswords(hashedPassword, password)
	}
}

// ============================================================================
// BENCHMARKS PARALELOS - Simula carga CONCURRENTE (múltiples usuarios)
// ============================================================================

// BenchmarkCreateUser_Parallel simula múltiples usuarios registrándose SIMULTÁNEAMENTE
// CRÍTICO: Este benchmark medirá el VERDADERO impacto de la concurrencia
// ANTES (sin optimizar): Bcrypt bloqueará, verás tiempos altos
// DESPUÉS (con goroutines): Deberías ver mejora si implementas worker pools
func BenchmarkCreateUser_Parallel(b *testing.B) {
	mockRepo := NewMockUserRepository()
	bcryptService := adapters.NewBcrypt()
	useCase := NewCreateUserUseCase(mockRepo, bcryptService)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			testUser := entities.User{
				Username:  "testuser",
				Email:     "test" + string(rune(counter)) + "@example.com",
				Password:  "Test123!@#",
				Nombre:    "Test",
				Apellidos: "User",
			}

			_, err := useCase.Execute(testUser)
			if err != nil && err.Error() != "el email test@example.com ya está registrado" {
				b.Fatalf("Error en CreateUser parallel: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkLogin_Parallel simula múltiples usuarios haciendo login SIMULTÁNEAMENTE
// CRÍTICO: Mide el throughput bajo carga concurrente
// ANTES: Verás que el sistema se satura rápidamente
// DESPUÉS: Con worker pools, deberías manejar más requests/segundo
func BenchmarkLogin_Parallel(b *testing.B) {
	mockRepo := NewMockUserRepository()
	bcryptService := adapters.NewBcrypt()
	mockJWT := &MockTokenManager{}
	useCase := NewLoginUseCase(mockRepo, mockJWT, bcryptService)

	hashedPassword, _ := bcryptService.HashPassword("Test123!@#")
	mockRepo.users["test@example.com"] = &entities.User{
		Id:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		Nombre:    "Test",
		Apellidos: "User",
	}

	input := LoginInput{
		Email:    "test@example.com",
		Password: "Test123!@#",
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := useCase.Execute(input)
			if err != nil {
				b.Fatalf("Error en Login parallel: %v", err)
			}
		}
	})
}

// ============================================================================
// BENCHMARKS DE CARGA PESADA - Para medir límites del sistema
// ============================================================================

// BenchmarkCreateUser_HighLoad simula carga extrema de creación de usuarios
// Ejecutar con: go test -bench=BenchmarkCreateUser_HighLoad -benchtime=100x
func BenchmarkCreateUser_HighLoad(b *testing.B) {
	mockRepo := NewMockUserRepository()
	bcryptService := adapters.NewBcrypt()
	useCase := NewCreateUserUseCase(mockRepo, bcryptService)

	b.SetParallelism(100) // Simula 100 goroutines concurrentes

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			testUser := entities.User{
				Username:  "loadtest",
				Email:     "load" + string(rune(counter)) + "@example.com",
				Password:  "Test123!@#",
				Nombre:    "Load",
				Apellidos: "Test",
			}

			_, err := useCase.Execute(testUser)
			if err != nil {
				b.Logf("Warning en high load: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkLogin_HighLoad simula carga extrema de logins
// Ejecutar con: go test -bench=BenchmarkLogin_HighLoad -benchtime=100x
func BenchmarkLogin_HighLoad(b *testing.B) {
	mockRepo := NewMockUserRepository()
	bcryptService := adapters.NewBcrypt()
	mockJWT := &MockTokenManager{}
	useCase := NewLoginUseCase(mockRepo, mockJWT, bcryptService)

	hashedPassword, _ := bcryptService.HashPassword("Test123!@#")
	mockRepo.users["test@example.com"] = &entities.User{
		Id:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		Nombre:    "Test",
		Apellidos: "User",
	}

	b.SetParallelism(100)

	input := LoginInput{
		Email:    "test@example.com",
		Password: "Test123!@#",
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := useCase.Execute(input)
			if err != nil {
				b.Logf("Warning en high load: %v", err)
			}
		}
	})
}

// ============================================================================
// BENCHMARKS COMPARATIVOS - Diferentes costs de Bcrypt
// ============================================================================

// BenchmarkBcrypt_Cost10 - Más rápido pero menos seguro
func BenchmarkBcrypt_Cost10(b *testing.B) {
	benchmarkBcryptCost(b, 10)
}

// BenchmarkBcrypt_Cost12 - Balance (tu configuración actual)
func BenchmarkBcrypt_Cost12(b *testing.B) {
	benchmarkBcryptCost(b, 12)
}

// BenchmarkBcrypt_Cost14 - Más seguro pero más lento
func BenchmarkBcrypt_Cost14(b *testing.B) {
	benchmarkBcryptCost(b, 14)
}

func benchmarkBcryptCost(b *testing.B, cost int) {
	password := "Test123!@#"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bcryptService := adapters.NewBcrypt()
		_, err := bcryptService.HashPassword(password)
		if err != nil {
			b.Fatalf("Error hashing with cost %d: %v", cost, err)
		}
	}
}