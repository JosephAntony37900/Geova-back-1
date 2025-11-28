// geova-back-1/Projects/application/projects_usecase_test.go
package application

import (
	"errors"
	"testing"
	"time"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
)

// Función de test dummy requerida para que Go reconozca el archivo
func TestDummy(t *testing.T) {
	// Esta función permite que go test ejecute los benchmarks
	t.Log("Benchmarks listos para ejecutar")
}
// ============================================================================
// MOCKS para testing de Projects
// ============================================================================

// MockProjectRepository simula el repositorio de proyectos
type MockProjectRepository struct {
	projects  map[int]*entities.Project
	saveDelay time.Duration
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects:  make(map[int]*entities.Project),
		saveDelay: 0,
	}
}

func (m *MockProjectRepository) Save(project entities.Project) error {
	if m.saveDelay > 0 {
		time.Sleep(m.saveDelay)
	}
	m.projects[project.Id] = &project
	return nil
}

func (m *MockProjectRepository) FindAll() ([]entities.Project, error) {
	projects := make([]entities.Project, 0)
	for _, p := range m.projects {
		projects = append(projects, *p)
	}
	return projects, nil
}

func (m *MockProjectRepository) FindById(id int) (*entities.Project, error) {
	if project, exists := m.projects[id]; exists {
		return project, nil
	}
	return nil, errors.New("proyecto no encontrado")
}

func (m *MockProjectRepository) FindByName(name string) ([]entities.Project, error) {
	projects := make([]entities.Project, 0)
	for _, p := range m.projects {
		if p.NombreProyecto == name {
			projects = append(projects, *p)
		}
	}
	return projects, nil
}

func (m *MockProjectRepository) FindByCategory(category string) ([]entities.Project, error) {
	projects := make([]entities.Project, 0)
	for _, p := range m.projects {
		if p.Categoria == category {
			projects = append(projects, *p)
		}
	}
	return projects, nil
}

func (m *MockProjectRepository) FindByDate(date string) ([]entities.Project, error) {
	projects := make([]entities.Project, 0)
	for _, p := range m.projects {
		if p.Fecha == date {
			projects = append(projects, *p)
		}
	}
	return projects, nil
}

func (m *MockProjectRepository) FindByUserId(userId int) ([]entities.Project, error) {
	projects := make([]entities.Project, 0)
	for _, p := range m.projects {
		if p.UserId == userId {
			projects = append(projects, *p)
		}
	}
	return projects, nil
}

func (m *MockProjectRepository) Update(project entities.Project) error {
	if m.saveDelay > 0 {
		time.Sleep(m.saveDelay)
	}
	m.projects[project.Id] = &project
	return nil
}

func (m *MockProjectRepository) Delete(id int) error {
	if _, exists := m.projects[id]; exists {
		delete(m.projects, id)
		return nil
	}
	return errors.New("proyecto no encontrado")
}

func (m *MockProjectRepository) GetProjectsStats(userId int, days int) ([]entities.DailyProjectCount, error) {
	return []entities.DailyProjectCount{}, nil
}

// ============================================================================
// MockCloudinaryService - Simula subidas de imágenes a Cloudinary
// ============================================================================

type MockCloudinaryService struct {
	uploadDelay time.Duration
	shouldFail  bool
}

func NewMockCloudinaryService() *MockCloudinaryService {
	return &MockCloudinaryService{
		uploadDelay: 0,
		shouldFail:  false,
	}
}

func (m *MockCloudinaryService) UploadImage(imagePath string) (string, error) {
	if m.uploadDelay > 0 {
		time.Sleep(m.uploadDelay)
	}
	if m.shouldFail {
		return "", errors.New("connection refused")
	}
	return "https://res.cloudinary.com/demo/image/upload/sample.jpg", nil
}

// ============================================================================
// Función helper para crear el worker service real con mock de Cloudinary
// ============================================================================

func createWorkerService(mockCloud *MockCloudinaryService) *services.ImageUploadWorkerService {
	// Crear el worker service REAL pero con el mock de Cloudinary
	// Esto es lo que el constructor espera: *services.ImageUploadWorkerService
	return services.NewImageUploadWorkerService(mockCloud, 2, 10)
}

// ============================================================================
// BENCHMARKS - CreateProject SIN IMAGEN (baseline)
// ============================================================================

func BenchmarkCreateProject(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockCloud := NewMockCloudinaryService()
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	testProject := entities.Project{
		Id:             1,
		NombreProyecto: "Benchmark Project",
		Fecha:          "2024-01-01",
		Categoria:      "desarrollo",
		Descripcion:    "Proyecto de prueba para benchmarks",
		Lat:            19.4326,
		Lng:            -99.1332,
		UserId:         1,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testProject.Id = i
		_, err := useCase.Execute(testProject, "")
		if err != nil {
			b.Fatalf("Error en CreateProject: %v", err)
		}
	}
}

// ============================================================================
// BENCHMARKS - CreateProject CON IMAGEN
// ============================================================================

func BenchmarkCreateProject_WithImage(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockCloud := NewMockCloudinaryService()
	mockCloud.uploadDelay = 100 * time.Millisecond
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	testProject := entities.Project{
		Id:             1,
		NombreProyecto: "Project With Image",
		Fecha:          "2024-01-01",
		Categoria:      "desarrollo",
		Descripcion:    "Proyecto con imagen",
		Lat:            19.4326,
		Lng:            -99.1332,
		UserId:         1,
	}
	imagePath := "/path/to/test/image.jpg"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testProject.Id = i
		_, err := useCase.Execute(testProject, imagePath)
		if err != nil {
			b.Fatalf("Error en CreateProject con imagen: %v", err)
		}
	}
}

func BenchmarkCreateProject_WithImage_FastUpload(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockCloud := NewMockCloudinaryService()
	mockCloud.uploadDelay = 50 * time.Millisecond
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	testProject := entities.Project{
		Id:             1,
		NombreProyecto: "Fast Upload Project",
		Fecha:          "2024-01-01",
		Categoria:      "desarrollo",
		Descripcion:    "Proyecto con upload rápido",
		Lat:            19.4326,
		Lng:            -99.1332,
		UserId:         1,
	}
	imagePath := "/path/to/test/image.jpg"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testProject.Id = i
		_, err := useCase.Execute(testProject, imagePath)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

func BenchmarkCreateProject_WithImage_SlowUpload(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockCloud := NewMockCloudinaryService()
	mockCloud.uploadDelay = 500 * time.Millisecond
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	testProject := entities.Project{
		Id:             1,
		NombreProyecto: "Slow Upload Project",
		Fecha:          "2024-01-01",
		Categoria:      "desarrollo",
		Descripcion:    "Proyecto con upload lento",
		Lat:            19.4326,
		Lng:            -99.1332,
		UserId:         1,
	}
	imagePath := "/path/to/test/image.jpg"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testProject.Id = i
		_, err := useCase.Execute(testProject, imagePath)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

func BenchmarkCreateProject_WithImage_Offline(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockCloud := NewMockCloudinaryService()
	mockCloud.shouldFail = true
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	testProject := entities.Project{
		Id:             1,
		NombreProyecto: "Offline Project",
		Fecha:          "2024-01-01",
		Categoria:      "desarrollo",
		Descripcion:    "Proyecto en modo offline",
		Lat:            19.4326,
		Lng:            -99.1332,
		UserId:         1,
	}
	imagePath := "/path/to/test/image.jpg"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testProject.Id = i
		result, err := useCase.Execute(testProject, imagePath)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
		if !result.IsOffline {
			b.Fatal("Expected offline mode")
		}
	}
}

// ============================================================================
// BENCHMARKS PARALELOS
// ============================================================================

func BenchmarkCreateProject_Parallel(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockCloud := NewMockCloudinaryService()
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			testProject := entities.Project{
				Id:             counter,
				NombreProyecto: "Parallel Project",
				Fecha:          "2024-01-01",
				Categoria:      "desarrollo",
				Descripcion:    "Proyecto paralelo",
				Lat:            19.4326,
				Lng:            -99.1332,
				UserId:         1,
			}

			_, err := useCase.Execute(testProject, "")
			if err != nil {
				b.Fatalf("Error en parallel: %v", err)
			}
			counter++
		}
	})
}

func BenchmarkCreateProject_WithImage_Parallel(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockCloud := NewMockCloudinaryService()
	mockCloud.uploadDelay = 100 * time.Millisecond
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			testProject := entities.Project{
				Id:             counter,
				NombreProyecto: "Parallel Image Project",
				Fecha:          "2024-01-01",
				Categoria:      "desarrollo",
				Descripcion:    "Proyecto con imagen paralelo",
				Lat:            19.4326,
				Lng:            -99.1332,
				UserId:         1,
			}
			imagePath := "/path/to/test/image.jpg"

			_, err := useCase.Execute(testProject, imagePath)
			if err != nil {
				b.Fatalf("Error en parallel con imagen: %v", err)
			}
			counter++
		}
	})
}

func BenchmarkCreateProject_HighLoad(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockCloud := NewMockCloudinaryService()
	mockCloud.uploadDelay = 100 * time.Millisecond
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	b.SetParallelism(100)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			testProject := entities.Project{
				Id:             counter,
				NombreProyecto: "High Load Project",
				Fecha:          "2024-01-01",
				Categoria:      "desarrollo",
				Descripcion:    "Proyecto bajo carga extrema",
				Lat:            19.4326,
				Lng:            -99.1332,
				UserId:         1,
			}
			imagePath := "/path/to/test/image.jpg"

			_, err := useCase.Execute(testProject, imagePath)
			if err != nil {
				b.Logf("Warning en high load: %v", err)
			}
			counter++
		}
	})
}

func BenchmarkCreateProject_VariableLatency(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockCloud := NewMockCloudinaryService()
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	testProject := entities.Project{
		Id:             1,
		NombreProyecto: "Variable Latency Project",
		Fecha:          "2024-01-01",
		Categoria:      "desarrollo",
		Descripcion:    "Proyecto con latencia variable",
		Lat:            19.4326,
		Lng:            -99.1332,
		UserId:         1,
	}
	imagePath := "/path/to/test/image.jpg"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		mockCloud.uploadDelay = time.Duration((i%5+1)*50) * time.Millisecond
		testProject.Id = i

		_, err := useCase.Execute(testProject, imagePath)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

func BenchmarkCreateProject_WithDBLatency(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockRepo.saveDelay = 50 * time.Millisecond
	mockCloud := NewMockCloudinaryService()
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	testProject := entities.Project{
		Id:             1,
		NombreProyecto: "DB Latency Project",
		Fecha:          "2024-01-01",
		Categoria:      "desarrollo",
		Descripcion:    "Proyecto con latencia en BD",
		Lat:            19.4326,
		Lng:            -99.1332,
		UserId:         1,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testProject.Id = i
		_, err := useCase.Execute(testProject, "")
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

func BenchmarkCreateProject_WithDBAndImageLatency(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockRepo.saveDelay = 50 * time.Millisecond
	mockCloud := NewMockCloudinaryService()
	mockCloud.uploadDelay = 200 * time.Millisecond
	workerSrv := createWorkerService(mockCloud)
	defer workerSrv.Shutdown()

	useCase := NewCreateProjectUseCase(mockRepo, mockCloud, workerSrv)

	testProject := entities.Project{
		Id:             1,
		NombreProyecto: "Combined Latency Project",
		Fecha:          "2024-01-01",
		Categoria:      "desarrollo",
		Descripcion:    "Proyecto con BD y upload lentos",
		Lat:            19.4326,
		Lng:            -99.1332,
		UserId:         1,
	}
	imagePath := "/path/to/test/image.jpg"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testProject.Id = i
		_, err := useCase.Execute(testProject, imagePath)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}