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
	Reply     chan ImageUploadResult // Reply channel for job-specific results
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
	shutdownMu  sync.Mutex // Protect shutdown state
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
	// Start worker goroutines
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i + 1)
	}

	// Start metrics processor (tracked by WaitGroup for safe shutdown)
	s.wg.Add(1)
	go s.processResults()
}

// worker es el proceso que consume trabajos del canal
func (s *ImageUploadWorkerService) worker(id int) {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			log.Printf("Worker #%d detenido", id)
			return
		case job, ok := <-s.jobQueue:
			if !ok {
				// Channel closed, exit worker
				log.Printf("Worker #%d detenido (canal cerrado)", id)
				return
			}
			s.processJob(id, job)
		}
	}
}

func (s *ImageUploadWorkerService) processJob(workerID int, job ImageUploadJob) {
	startTime := time.Now()

	if job.ProjectID <= 0 {
		result := ImageUploadResult{
			ProjectID: job.ProjectID,
			Error:     fmt.Errorf("ID inválido"),
			Duration:  time.Since(startTime),
		}
		s.sendResult(job, result)
		return
	}

	// Retry logic with exponential backoff
	maxAttempts := 3
	delays := []time.Duration{500 * time.Millisecond, 1 * time.Second, 2 * time.Second}

	var url string
	var uploadErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			delay := delays[attempt-1]
			log.Printf("Worker #%d - Reintentando subida (intento %d/%d) para proyecto ID=%d después de %v",
				workerID, attempt+1, maxAttempts, job.ProjectID, delay)
			time.Sleep(delay)
		}

		url, uploadErr = s.cloudSrv.UploadImage(job.ImagePath)
		if uploadErr == nil {
			// Success!
			break
		}

		log.Printf("Worker #%d - Intento %d/%d fallido para proyecto ID=%d: %v",
			workerID, attempt+1, maxAttempts, job.ProjectID, uploadErr)
	}

	if uploadErr != nil {
		log.Printf("Worker #%d - Error definitivo al subir imagen después de %d intentos: %v",
			workerID, maxAttempts, uploadErr)
		result := ImageUploadResult{
			ProjectID: job.ProjectID,
			Error:     uploadErr,
			Duration:  time.Since(startTime),
		}
		s.sendResult(job, result)
		return
	}

	// Image uploaded successfully, now update database
	project, err := s.repo.FindById(job.ProjectID)
	if err != nil {
		log.Printf("Worker #%d - No se encontró proyecto ID=%d: %v", workerID, job.ProjectID, err)
		result := ImageUploadResult{
			ProjectID: job.ProjectID,
			Error:     err,
			Duration:  time.Since(startTime),
		}
		s.sendResult(job, result)
		return
	}

	// Update image URL
	project.Img = url

	// Protect concurrent writes with mutex
	s.mu.Lock()
	err = s.repo.Update(*project)
	s.mu.Unlock()

	duration := time.Since(startTime)

	if err != nil {
		log.Printf("Worker #%d - Error al actualizar proyecto: %v", workerID, err)
		result := ImageUploadResult{
			ProjectID: job.ProjectID,
			Error:     err,
			Duration:  duration,
		}
		s.sendResult(job, result)
		return
	}

	log.Printf("Worker #%d - Proyecto ID=%d actualizado exitosamente (duración: %v)",
		workerID, job.ProjectID, duration)

	result := ImageUploadResult{
		ProjectID: job.ProjectID,
		ImageURL:  url,
		Error:     nil,
		Duration:  duration,
	}
	s.sendResult(job, result)
}

// sendResult sends result to both reply channel (if present) and resultQueue for metrics
func (s *ImageUploadWorkerService) sendResult(job ImageUploadJob, result ImageUploadResult) {
	// Send to job-specific reply channel if present (non-blocking)
	if job.Reply != nil {
		select {
		case job.Reply <- result:
			// Successfully sent to reply channel
		default:
			log.Printf("WARNING: reply channel full or not listened for project ID=%d", job.ProjectID)
		}
	}

	// Send to resultQueue for metrics (non-blocking)
	select {
	case s.resultQueue <- result:
		// Successfully sent to metrics queue
	default:
		// Queue full, skip metrics for this result
		log.Printf("WARNING: resultQueue full, skipping metrics for project ID=%d", job.ProjectID)
	}
}

func (s *ImageUploadWorkerService) processResults() {
	defer s.wg.Done()

	for result := range s.resultQueue {
		if result.Error != nil {
			log.Printf("METRICS: Subida fallida - Proyecto=%d, Error=%v, Duración=%v",
				result.ProjectID, result.Error, result.Duration)
		} else {
			log.Printf("METRICS: Subida exitosa - Proyecto=%d, Duración=%v",
				result.ProjectID, result.Duration)
		}
	}
	log.Println("METRICS: processResults finalizado")
}

func (s *ImageUploadWorkerService) SubmitUploadJob(projectID int, imagePath string) error {
	// Check if service is shutting down
	select {
	case <-s.shutdown:
		return fmt.Errorf("service shutting down")
	default:
	}

	job := ImageUploadJob{
		ProjectID: projectID,
		ImagePath: imagePath,
		Timestamp: time.Now(),
		Reply:     nil, // No reply channel for async jobs
	}

	select {
	case s.jobQueue <- job:
		log.Printf("Trabajo encolado - Proyecto ID=%d", projectID)
		return nil
	case <-s.shutdown:
		return fmt.Errorf("service shutting down")
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout: cola de trabajos llena")
	}
}

// SubmitUploadJobSync envía un trabajo y espera el resultado usando un reply channel
func (s *ImageUploadWorkerService) SubmitUploadJobSync(projectID int, imagePath string, timeout time.Duration) (*ImageUploadResult, error) {
	// Check if service is shutting down
	select {
	case <-s.shutdown:
		return nil, fmt.Errorf("service shutting down")
	default:
	}

	// Create buffered reply channel
	reply := make(chan ImageUploadResult, 1)

	job := ImageUploadJob{
		ProjectID: projectID,
		ImagePath: imagePath,
		Timestamp: time.Now(),
		Reply:     reply,
	}

	// Enqueue job
	select {
	case s.jobQueue <- job:
		log.Printf("Trabajo síncrono encolado - Proyecto ID=%d", projectID)
	case <-s.shutdown:
		return nil, fmt.Errorf("service shutting down")
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout: cola de trabajos llena")
	}

	// Wait for result on reply channel
	select {
	case result := <-reply:
		return &result, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout esperando resultado después de %v", timeout)
	}
}

// Shutdown detiene gracefully todos los workers
func (s *ImageUploadWorkerService) Shutdown() {
	s.shutdownMu.Lock()
	defer s.shutdownMu.Unlock()

	// Check if already shut down
	select {
	case <-s.shutdown:
		// Already shut down
		return
	default:
	}

	log.Println("Iniciando shutdown del servicio de imágenes...")

	// 1. Signal shutdown to reject new jobs
	close(s.shutdown)

	// 2. Close job queue so workers finish processing
	close(s.jobQueue)

	// 3. Wait for all workers to complete
	s.wg.Wait()
	log.Println("Todos los workers finalizados")

	// 4. Close result queue after workers are done
	close(s.resultQueue)

	log.Println("Servicio de imágenes detenido correctamente")
}

func (s *ImageUploadWorkerService) GetQueueSize() int {
	return len(s.jobQueue)
}
