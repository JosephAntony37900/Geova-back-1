// geova-back-1/Projects/application/projects_usecase_test.go
package application

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/entities"
	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/services"
)

// ============================================================================
// TEST HELPERS Y UTILITIES
// ============================================================================

// Función de test dummy requerida para que Go reconozca el archivo
func TestDummy(t *testing.T) {
	t.Log("Benchmarks listos para ejecutar")
	printBenchmarkInstructions(t)
}

func printBenchmarkInstructions(t *testing.T) {
	t.Log("\n" + string([]byte{27}) + "[1;36m" + "═══════════════════════════════════════════════════════════════")
	t.Log("  GUÍA DE EJECUCIÓN DE BENCHMARKS")
	t.Log("═══════════════════════════════════════════════════════════════" + string([]byte{27}) + "[0m")
	t.Log("\n📊 Para ejecutar TODOS los benchmarks:")
	t.Log("   go test -bench=. -benchmem -benchtime=3s")
	t.Log("\n🎯 Para ejecutar benchmarks específicos:")
	t.Log("   go test -bench=BenchmarkCreateProject -benchmem")
	t.Log("   go test -bench=BenchmarkGetTotalProjects -benchmem")
	t.Log("   go test -bench=BenchmarkGetProjectStats -benchmem")
	t.Log("\n⚡ Para benchmarks paralelos con alta carga:")
	t.Log("   go test -bench=Parallel -benchmem -cpu=1,2,4,8")
	t.Log("\n💾 Para guardar resultados:")
	t.Log("   go test -bench=. -benchmem > benchmark_results.txt")
	t.Log("\n🔍 Para comparar con resultados anteriores:")
	t.Log("   go test -bench=. -benchmem > new.txt")
	t.Log("   benchstat old.txt new.txt")
	t.Log("")
}

// BenchmarkResult ayuda a formatear y mostrar resultados
type BenchmarkResult struct {
	Name           string
	Operations     int
	NsPerOp        int64
	BytesPerOp     int64
	AllocsPerOp    int64
	ParallelProcs  int
}

func printBenchmarkSummary(b *testing.B, category string) {
	b.Log("\n" + string([]byte{27}) + "[1;32m" + "✓ Benchmark completado: " + category + string([]byte{27}) + "[0m")
}

// ============================================================================
// MOCKS MEJORADOS para testing de Projects
// ============================================================================

// MockProjectRepository simula el repositorio de proyectos
type MockProjectRepository struct {
	projects  map[int]*entities.Project
	saveDelay time.Duration
	findDelay time.Duration
	totalCallCount int
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects:  make(map[int]*entities.Project),
		saveDelay: 0,
		findDelay: 0,
		totalCallCount: 0,
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
	if m.findDelay > 0 {
		time.Sleep(m.findDelay)
	}
	projects := make([]entities.Project, 0)
	for _, p := range m.projects {
		projects = append(projects, *p)
	}
	return projects, nil
}

func (m *MockProjectRepository) FindById(id int) (*entities.Project, error) {
	if m.findDelay > 0 {
		time.Sleep(m.findDelay)
	}
	if project, exists := m.projects[id]; exists {
		return project, nil
	}
	return nil, errors.New("proyecto no encontrado")
}

func (m *MockProjectRepository) FindByName(name string) ([]entities.Project, error) {
	if m.findDelay > 0 {
		time.Sleep(m.findDelay)
	}
	projects := make([]entities.Project, 0)
	for _, p := range m.projects {
		if p.NombreProyecto == name {
			projects = append(projects, *p)
		}
	}
	return projects, nil
}

func (m *MockProjectRepository) FindByCategory(category string) ([]entities.Project, error) {
	if m.findDelay > 0 {
		time.Sleep(m.findDelay)
	}
	projects := make([]entities.Project, 0)
	for _, p := range m.projects {
		if p.Categoria == category {
			projects = append(projects, *p)
		}
	}
	return projects, nil
}

func (m *MockProjectRepository) FindByDate(date string) ([]entities.Project, error) {
	if m.findDelay > 0 {
		time.Sleep(m.findDelay)
	}
	projects := make([]entities.Project, 0)
	for _, p := range m.projects {
		if p.Fecha == date {
			projects = append(projects, *p)
		}
	}
	return projects, nil
}

func (m *MockProjectRepository) FindByUserId(userId int) ([]entities.Project, error) {
	if m.findDelay > 0 {
		time.Sleep(m.findDelay)
	}
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
	if m.findDelay > 0 {
		time.Sleep(m.findDelay)
	}
	return []entities.DailyProjectCount{}, nil
}

// NUEVO: Método para obtener total de proyectos
func (m *MockProjectRepository) GetTotalProjectsByUser(userId string) (int, error) {
	m.totalCallCount++
	if m.findDelay > 0 {
		time.Sleep(m.findDelay)
	}
	
	count := 0
	for _, p := range m.projects {
		if fmt.Sprintf("%d", p.UserId) == userId {
			count++
		}
	}
	return count, nil
}

// Helper para poblar datos de prueba
func (m *MockProjectRepository) SeedTestData(userId int, count int) {
	for i := 0; i < count; i++ {
		project := entities.Project{
			Id:             i + 1,
			NombreProyecto: fmt.Sprintf("Test Project %d", i+1),
			Fecha:          "2024-01-01",
			Categoria:      "desarrollo",
			Descripcion:    "Proyecto de prueba",
			Lat:            19.4326,
			Lng:            -99.1332,
			UserId:         userId,
		}
		m.projects[project.Id] = &project
	}
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
	return services.NewImageUploadWorkerService(mockCloud, 2, 10)
}

// ============================================================================
// BENCHMARKS - CreateProject (Existentes mejorados)
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
	
	printBenchmarkSummary(b, "CreateProject (sin imagen)")
}

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
	
	printBenchmarkSummary(b, "CreateProject (con imagen 100ms)")
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
	
	printBenchmarkSummary(b, "CreateProject (upload rápido 50ms)")
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
	
	printBenchmarkSummary(b, "CreateProject (upload lento 500ms)")
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
	
	printBenchmarkSummary(b, "CreateProject (modo offline)")
}

// ============================================================================
// NUEVOS BENCHMARKS - GetTotalProjects
// ============================================================================

func BenchmarkGetTotalProjects_Empty(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	useCase := NewGetTotalProjectsByUserUseCase(mockRepo)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := useCase.Execute("1")
		if err != nil {
			b.Fatalf("Error en GetTotalProjects: %v", err)
		}
	}
	
	printBenchmarkSummary(b, "GetTotalProjects (sin proyectos)")
}

func BenchmarkGetTotalProjects_Small(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockRepo.SeedTestData(1, 10) // 10 proyectos
	useCase := NewGetTotalProjectsByUserUseCase(mockRepo)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		count, err := useCase.Execute("1")
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
		if count != 10 {
			b.Fatalf("Expected 10 projects, got %d", count)
		}
	}
	
	printBenchmarkSummary(b, "GetTotalProjects (10 proyectos)")
}

func BenchmarkGetTotalProjects_Medium(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockRepo.SeedTestData(1, 100) // 100 proyectos
	useCase := NewGetTotalProjectsByUserUseCase(mockRepo)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		count, err := useCase.Execute("1")
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
		if count != 100 {
			b.Fatalf("Expected 100 projects, got %d", count)
		}
	}
	
	printBenchmarkSummary(b, "GetTotalProjects (100 proyectos)")
}

func BenchmarkGetTotalProjects_Large(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockRepo.SeedTestData(1, 1000) // 1000 proyectos
	useCase := NewGetTotalProjectsByUserUseCase(mockRepo)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		count, err := useCase.Execute("1")
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
		if count != 1000 {
			b.Fatalf("Expected 1000 projects, got %d", count)
		}
	}
	
	printBenchmarkSummary(b, "GetTotalProjects (1000 proyectos)")
}

func BenchmarkGetTotalProjects_WithDBLatency(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockRepo.SeedTestData(1, 100)
	mockRepo.findDelay = 50 * time.Millisecond
	useCase := NewGetTotalProjectsByUserUseCase(mockRepo)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := useCase.Execute("1")
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
	
	printBenchmarkSummary(b, "GetTotalProjects (con latencia BD 50ms)")
}

func BenchmarkGetTotalProjects_Parallel(b *testing.B) {
	mockRepo := NewMockProjectRepository()
	mockRepo.SeedTestData(1, 100)
	useCase := NewGetTotalProjectsByUserUseCase(mockRepo)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := useCase.Execute("1")
			if err != nil {
				b.Fatalf("Error en parallel: %v", err)
			}
		}
	})
	
	printBenchmarkSummary(b, "GetTotalProjects (paralelo)")
}



// ============================================================================
// BENCHMARKS PARALELOS EXISTENTES (Mejorados)
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
	
	printBenchmarkSummary(b, "CreateProject (paralelo sin imagen)")
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
	
	printBenchmarkSummary(b, "CreateProject (paralelo con imagen)")
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
	
	printBenchmarkSummary(b, "CreateProject (carga alta)")
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

printBenchmarkSummary(b, "CreateProject (latencia variable)")
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

printBenchmarkSummary(b, "CreateProject (latencia BD 50ms)")
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

printBenchmarkSummary(b, "CreateProject (latencia combinada BD+Upload)")
}
