package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/JosephAntony37900/Geova-back-1/Projects/domain/repository"
)

type ImageUploadJob struct {
	ProjectID int
	ImagePath string
	Timestamp time.Time
}

type ImageUploadResult struct {
	ProjectID int
	ImageURL  string
	Error     error
	Duration  time.Duration
}

type ImageUploadWorkerService struct {
	cloudSrv    ICloudinaryService
	repo        repository.ProjectRepository
	jobQueue    chan ImageUploadJob
	resultQueue chan ImageUploadResult
	workers     int
	wg          sync.WaitGroup
	mu          sync.Mutex
	shutdown    chan struct{}
}

// NewImageUploadWorkerService crea el servicio con workers
func NewImageUploadWorkerService(
	cloudSrv ICloudinaryService,
	repo repository.ProjectRepository,
	numWorkers int,
) *ImageUploadWorkerService {
	if numWorkers <= 0 {
		numWorkers = 3 
	}

	service := &ImageUploadWorkerService{
		cloudSrv:    cloudSrv,
		repo:        repo,
		jobQueue:    make(chan ImageUploadJob, 100),
		resultQueue: make(chan ImageUploadResult, 100),
		workers:     numWorkers,
		shutdown:    make(chan struct{}),
	}

	service.Start()

	return service
}

func (s *ImageUploadWorkerService) Start() {
	
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i + 1)
	}

	go s.processResults()
}

// worker es el proceso que consume trabajos del canal
func (s *ImageUploadWorkerService) worker(id int) {
	defer s.wg.Done()
	
	// log.Printf(" Worker #%d iniciado", id)

	for {
		select {
		case job := <-s.jobQueue:
			s.processJob(id, job)
		case <-s.shutdown:
			log.Printf("Worker #%d detenido", id)
			return
		}
	}
}

func (s *ImageUploadWorkerService) processJob(workerID int, job ImageUploadJob) {
	startTime := time.Now()
	
	//log.Printf("Worker #%d procesando proyecto ID=%d", workerID, job.ProjectID)

	if job.ProjectID <= 0 {
		//log.Printf("Worker #%d - ID de proyecto inválido: %d", workerID, job.ProjectID)
		s.resultQueue <- ImageUploadResult{
			ProjectID: job.ProjectID,
			Error:     fmt.Errorf("ID inválido"),
		}
		return
	}

	// Subir imagen a Cloudinary
	url, err := s.cloudSrv.UploadImage(job.ImagePath)
	if err != nil {
		log.Printf("Worker #%d - Error al subir imagen: %v", workerID, err)
		s.resultQueue <- ImageUploadResult{
			ProjectID: job.ProjectID,
			Error:     err,
			Duration:  time.Since(startTime),
		}
		return
	}

	//log.Printf(" Worker #%d - Imagen subida: %s", workerID, url)

	// Buscar proyecto en BD
	project, err := s.repo.FindById(job.ProjectID)
	if err != nil {
		log.Printf("❌ Worker #%d - No se encontró proyecto ID=%d: %v", workerID, job.ProjectID, err)
		s.resultQueue <- ImageUploadResult{
			ProjectID: job.ProjectID,
			Error:     err,
			Duration:  time.Since(startTime),
		}
		return
	}

	// Actualizar imagen
	project.Img = url

	// Proteger escritura concurrente con mutex
	s.mu.Lock()
	err = s.repo.Update(project)
	s.mu.Unlock()

	if err != nil {
		log.Printf("Worker #%d - Error al actualizar proyecto: %v", workerID, err)
		s.resultQueue <- ImageUploadResult{
			ProjectID: job.ProjectID,
			Error:     err,
			Duration:  time.Since(startTime),
		}
		return
	}

	duration := time.Since(startTime)
	log.Printf(" Worker #%d - Proyecto ID=%d actualizado (duración: %v)", 
		workerID, job.ProjectID, duration)

	s.resultQueue <- ImageUploadResult{
		ProjectID: job.ProjectID,
		ImageURL:  url,
		Error:     nil,
		Duration:  duration,
	}
}

func (s *ImageUploadWorkerService) processResults() {
	for result := range s.resultQueue {
		if result.Error != nil {
			log.Printf("METRICS: Subida fallida - Proyecto=%d, Error=%v, Duración=%v",
				result.ProjectID, result.Error, result.Duration)
		} else {
			log.Printf("METRICS: Subida exitosa - Proyecto=%d, Duración=%v",
				result.ProjectID, result.Duration)
		}
	}
}

func (s *ImageUploadWorkerService) SubmitUploadJob(projectID int, imagePath string) error {
	job := ImageUploadJob{
		ProjectID: projectID,
		ImagePath: imagePath,
		Timestamp: time.Now(),
	}

	select {
	case s.jobQueue <- job:
		log.Printf("Trabajo encolado - Proyecto ID=%d", projectID)
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout: cola de trabajos llena")
	}
}

// SubmitUploadJobSync envía un trabajo y espera el resultado 
func (s *ImageUploadWorkerService) SubmitUploadJobSync(projectID int, imagePath string, timeout time.Duration) (*ImageUploadResult, error) {
	job := ImageUploadJob{
		ProjectID: projectID,
		ImagePath: imagePath,
		Timestamp: time.Now(),
	}

	// Enviar trabajo
	select {
	case s.jobQueue <- job:
		log.Printf("Trabajo síncrono encolado - Proyecto ID=%d", projectID)
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout: cola de trabajos llena")
	}

	// Esperar resultado específico de este proyecto
	deadline := time.After(timeout)
	for {
		select {
		case result := <-s.resultQueue:
			if result.ProjectID == projectID {
				return &result, nil
			}
			// Si no es el nuestro, devolverlo a la cola
			go func(r ImageUploadResult) {
				s.resultQueue <- r
			}(result)
		case <-deadline:
			return nil, fmt.Errorf("timeout esperando resultado después de %v", timeout)
		}
	}
}

// Shutdown detiene gracefully todos los workers
func (s *ImageUploadWorkerService) Shutdown() {
	//log.Println("Iniciando shutdown del servicio de imágenes...")
	
	close(s.shutdown)
	close(s.jobQueue)
	
	s.wg.Wait()
	
	close(s.resultQueue)
	
	//log.Println("Servicio de imágenes detenido correctamente")
}

func (s *ImageUploadWorkerService) GetQueueSize() int {
	return len(s.jobQueue)
}