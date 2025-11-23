package services

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// ImageUploadJob representa un trabajo de subida de imagen con canal de respuesta
type ImageUploadJob struct {
	ProjectID int
	ImagePath string
	Timestamp time.Time
	Reply     chan ImageUploadResult // Canal para respuesta exclusiva del caller
}

// ImageUploadResult representa el resultado de una subida de imagen
type ImageUploadResult struct {
	ProjectID int
	ImageURL  string
	Error     error
	Duration  time.Duration
	Timestamp time.Time
}

// ImageUploadWorkerService gestiona un pool de workers para subir imágenes
type ImageUploadWorkerService struct {
	cloudSrv    ICloudinaryService
	jobQueue    chan ImageUploadJob
	resultQueue chan ImageUploadResult
	shutdown    chan struct{}
	wg          sync.WaitGroup
	mu          sync.Mutex // Protege acceso a repo si no es thread-safe
	numWorkers  int
}

// NewImageUploadWorkerService crea una nueva instancia del servicio
func NewImageUploadWorkerService(cloudSrv ICloudinaryService, numWorkers int, queueSize int) *ImageUploadWorkerService {
	if numWorkers <= 0 {
		numWorkers = 3
	}
	if queueSize <= 0 {
		queueSize = 100
	}

	service := &ImageUploadWorkerService{
		cloudSrv:    cloudSrv,
		jobQueue:    make(chan ImageUploadJob, queueSize),
		resultQueue: make(chan ImageUploadResult, queueSize),
		shutdown:    make(chan struct{}),
		numWorkers:  numWorkers,
	}

	// Iniciar workers
	for i := 1; i <= numWorkers; i++ {
		service.wg.Add(1)
		go service.worker(i)
	}

	// Iniciar procesador de resultados para métricas
	go service.processResults()

	log.Printf("ImageUploadWorkerService iniciado con %d workers y cola de %d", numWorkers, queueSize)
	return service
}

// worker es la goroutine que procesa trabajos
func (s *ImageUploadWorkerService) worker(id int) {
	defer s.wg.Done()
	log.Printf("Worker #%d: iniciado", id)

	for {
		select {
		case job, ok := <-s.jobQueue:
			if !ok {
				log.Printf("Worker #%d: jobQueue cerrado, saliendo", id)
				return
			}
			s.processJob(id, job)
		case <-s.shutdown:
			log.Printf("Worker #%d: detenido por shutdown", id)
			return
		}
	}
}

// processJob procesa un trabajo individual
func (s *ImageUploadWorkerService) processJob(workerID int, job ImageUploadJob) {
	startTime := time.Now()
	log.Printf("Worker #%d: procesando proyecto %d, imagen: %s", workerID, job.ProjectID, job.ImagePath)

	// Subir imagen a Cloudinary
	imageURL, err := s.cloudSrv.UploadImage(job.ImagePath)
	duration := time.Since(startTime)

	// Construir resultado
	result := ImageUploadResult{
		ProjectID: job.ProjectID,
		ImageURL:  imageURL,
		Error:     err,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		log.Printf("Worker #%d: error subiendo imagen proyecto %d: %v (duración: %v)",
			workerID, job.ProjectID, err, duration)
	} else {
		log.Printf("Worker #%d: imagen subida exitosamente proyecto %d: %s (duración: %v)",
			workerID, job.ProjectID, imageURL, duration)
	}

	// Enviar resultado al canal Reply si existe (para SubmitUploadJobSync)
	if job.Reply != nil {
		select {
		case job.Reply <- result:
			// Enviado exitosamente
		default:
			log.Printf("Advertencia: reply canal lleno o cerrado para proyecto %d", job.ProjectID)
		}
	}

	// Enviar copia a resultQueue para métricas (si está abierto)
	select {
	case s.resultQueue <- result:
		// Enviado exitosamente
	default:
		// Canal lleno o cerrado, no bloquear
	}
}

// processResults procesa resultados para métricas y logging
func (s *ImageUploadWorkerService) processResults() {
	log.Println("processResults: iniciado")
	successCount := 0
	errorCount := 0

	for result := range s.resultQueue {
		if result.Error != nil {
			errorCount++
			log.Printf("Métrica: Error en proyecto %d - Total errores: %d", result.ProjectID, errorCount)
		} else {
			successCount++
			log.Printf("Métrica: Éxito en proyecto %d - Total éxitos: %d, duración: %v",
				result.ProjectID, successCount, result.Duration)
		}
	}

	log.Printf("processResults: finalizado - Éxitos: %d, Errores: %d", successCount, errorCount)
}

// SubmitUploadJob encola un trabajo de forma asíncrona (sin esperar resultado)
func (s *ImageUploadWorkerService) SubmitUploadJob(projectID int, imagePath string) error {
	// Verificar si el servicio está en shutdown
	select {
	case <-s.shutdown:
		return fmt.Errorf("servicio en proceso de cierre, no se aceptan nuevos trabajos")
	default:
	}

	job := ImageUploadJob{
		ProjectID: projectID,
		ImagePath: imagePath,
		Timestamp: time.Now(),
		Reply:     nil, // No necesita respuesta directa
	}

	// Intentar encolar con timeout
	select {
	case s.jobQueue <- job:
		log.Printf("Trabajo encolado (async) para proyecto %d", projectID)
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout al encolar trabajo para proyecto %d", projectID)
	case <-s.shutdown:
		return fmt.Errorf("servicio cerrado durante encolamiento")
	}
}

// SubmitUploadJobSync encola un trabajo y espera el resultado de forma síncrona
func (s *ImageUploadWorkerService) SubmitUploadJobSync(projectID int, imagePath string, timeout time.Duration) (*ImageUploadResult, error) {
	// Verificar si el servicio está en shutdown
	select {
	case <-s.shutdown:
		return nil, fmt.Errorf("servicio en proceso de cierre, no se aceptan nuevos trabajos")
	default:
	}

	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	// Crear canal de respuesta con buffer de 1
	reply := make(chan ImageUploadResult, 1)

	job := ImageUploadJob{
		ProjectID: projectID,
		ImagePath: imagePath,
		Timestamp: time.Now(),
		Reply:     reply,
	}

	// Intentar encolar con timeout
	select {
	case s.jobQueue <- job:
		log.Printf("Trabajo encolado (sync) para proyecto %d", projectID)
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout al encolar trabajo para proyecto %d", projectID)
	case <-s.shutdown:
		return nil, fmt.Errorf("servicio cerrado durante encolamiento")
	}

	// Esperar resultado en el canal Reply o timeout
	select {
	case result := <-reply:
		if result.Error != nil {
			return &result, fmt.Errorf("error al subir imagen: %w", result.Error)
		}
		return &result, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout esperando resultado de subida para proyecto %d", projectID)
	case <-s.shutdown:
		return nil, fmt.Errorf("servicio cerrado mientras esperaba resultado")
	}
}

// Shutdown detiene el servicio de forma ordenada
func (s *ImageUploadWorkerService) Shutdown() {
	log.Println("ImageUploadWorkerService: iniciando shutdown...")

	// Cerrar canal de shutdown para señalar a workers y métodos
	close(s.shutdown)

	// Cerrar jobQueue para que workers terminen después de procesar trabajos pendientes
	close(s.jobQueue)

	// Esperar a que todos los workers terminen
	s.wg.Wait()
	log.Println("ImageUploadWorkerService: todos los workers han terminado")

	// Ahora es seguro cerrar resultQueue
	close(s.resultQueue)

	log.Println("ImageUploadWorkerService: shutdown completado")
}
