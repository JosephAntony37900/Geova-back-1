package services

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// ImageUploadJob representa un trabajo de subida de imagen
type ImageUploadJob struct {
	LocalPath string
	Reply     chan ImageUploadResult // Canal de respuesta para este job específico
}

// ImageUploadResult contiene el resultado de una subida de imagen
type ImageUploadResult struct {
	URL   string
	Error error
}

// ImageUploadWorkerService maneja la subida asíncrona de imágenes con workers
type ImageUploadWorkerService struct {
	cloudSrv    ICloudinaryService
	jobQueue    chan ImageUploadJob
	resultQueue chan ImageUploadResult
	shutdown    chan struct{}
	wg          sync.WaitGroup
	mu          sync.RWMutex // Protege el estado de shutdown
	isShutdown  bool
}

// NewImageUploadWorkerService crea un nuevo servicio de workers para subida de imágenes
func NewImageUploadWorkerService(cloudSrv ICloudinaryService, numWorkers int, queueSize int) *ImageUploadWorkerService {
	service := &ImageUploadWorkerService{
		cloudSrv:    cloudSrv,
		jobQueue:    make(chan ImageUploadJob, queueSize),
		resultQueue: make(chan ImageUploadResult, queueSize),
		shutdown:    make(chan struct{}),
		isShutdown:  false,
	}

	// Iniciar workers
	for i := 0; i < numWorkers; i++ {
		service.wg.Add(1)
		go service.worker(i)
	}

	log.Printf("INFO: ImageUploadWorkerService iniciado con %d workers", numWorkers)
	return service
}

// worker es el loop principal de cada worker
func (s *ImageUploadWorkerService) worker(id int) {
	defer s.wg.Done()
	log.Printf("INFO: Worker %d iniciado", id)

	for {
		select {
		case job, ok := <-s.jobQueue:
			if !ok {
				// jobQueue cerrado, terminar worker
				log.Printf("INFO: Worker %d terminando (jobQueue cerrado)", id)
				return
			}
			s.processJob(id, job)
		case <-s.shutdown:
			// Señal de shutdown recibida
			log.Printf("INFO: Worker %d recibió señal de shutdown", id)
			return
		}
	}
}

// processJob procesa un trabajo individual con reintentos
func (s *ImageUploadWorkerService) processJob(workerID int, job ImageUploadJob) {
	maxAttempts := 3
	delays := []time.Duration{500 * time.Millisecond, 1 * time.Second, 2 * time.Second}

	var url string
	var err error

	// Reintentar con backoff
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			delay := delays[attempt-1]
			log.Printf("INFO: Worker %d reintentando subida (intento %d/%d) tras %v", workerID, attempt+1, maxAttempts, delay)
			time.Sleep(delay)
		}

		url, err = s.cloudSrv.UploadImage(job.LocalPath)
		if err == nil {
			// Éxito
			log.Printf("SUCCESS: Worker %d subió imagen exitosamente: %s", workerID, url)
			break
		}

		log.Printf("WARNING: Worker %d falló intento %d/%d para subir imagen: %v", workerID, attempt+1, maxAttempts, err)
	}

	// Construir resultado
	result := ImageUploadResult{
		URL:   url,
		Error: err,
	}

	// Enviar resultado al canal Reply del job (no bloqueante)
	if job.Reply != nil {
		select {
		case job.Reply <- result:
			log.Printf("DEBUG: Worker %d envió resultado a Reply channel", workerID)
		default:
			log.Printf("WARNING: Worker %d no pudo enviar resultado a Reply channel (canal lleno o cerrado)", workerID)
		}
	}

	// Enviar copia a resultQueue para métricas (no bloqueante)
	select {
	case s.resultQueue <- result:
		log.Printf("DEBUG: Worker %d envió resultado a resultQueue", workerID)
	default:
		log.Printf("WARNING: Worker %d no pudo enviar resultado a resultQueue (cola llena)", workerID)
	}
}

// SubmitUploadJob encola un trabajo de forma asíncrona
func (s *ImageUploadWorkerService) SubmitUploadJob(localPath string) error {
	s.mu.RLock()
	if s.isShutdown {
		s.mu.RUnlock()
		return fmt.Errorf("servicio en shutdown, no se aceptan nuevos trabajos")
	}

	job := ImageUploadJob{
		LocalPath: localPath,
		Reply:     nil, // Sin canal de respuesta para trabajos asíncronos
	}

	select {
	case s.jobQueue <- job:
		s.mu.RUnlock()
		log.Printf("INFO: Trabajo encolado para: %s", localPath)
		return nil
	default:
		s.mu.RUnlock()
		return fmt.Errorf("cola de trabajos llena, intente más tarde")
	}
}

// SubmitUploadJobSync encola un trabajo y espera el resultado de forma síncrona
func (s *ImageUploadWorkerService) SubmitUploadJobSync(localPath string, timeout time.Duration) (string, error) {
	s.mu.RLock()
	if s.isShutdown {
		s.mu.RUnlock()
		return "", fmt.Errorf("servicio en shutdown, no se aceptan nuevos trabajos")
	}

	// Crear canal de respuesta con buffer
	reply := make(chan ImageUploadResult, 1)

	job := ImageUploadJob{
		LocalPath: localPath,
		Reply:     reply,
	}

	// Encolar el trabajo mientras mantenemos el lock
	select {
	case s.jobQueue <- job:
		s.mu.RUnlock()
		log.Printf("INFO: Trabajo síncrono encolado para: %s", localPath)
	default:
		s.mu.RUnlock()
		return "", fmt.Errorf("cola de trabajos llena, intente más tarde")
	}

	// Esperar resultado o timeout
	select {
	case result := <-reply:
		if result.Error != nil {
			return "", result.Error
		}
		return result.URL, nil
	case <-time.After(timeout):
		// En caso de timeout, el resultado podría aún llegar al canal buffered
		// pero no causará leak porque el buffer es de tamaño 1
		log.Printf("WARNING: Timeout esperando resultado de subida para: %s", localPath)
		return "", fmt.Errorf("timeout esperando resultado de subida de imagen")
	}
}

// GetResultQueue retorna el canal de resultados para métricas
func (s *ImageUploadWorkerService) GetResultQueue() <-chan ImageUploadResult {
	return s.resultQueue
}

// Shutdown cierra el servicio de forma ordenada
func (s *ImageUploadWorkerService) Shutdown() {
	s.mu.Lock()
	if s.isShutdown {
		s.mu.Unlock()
		log.Println("INFO: Shutdown ya fue llamado anteriormente")
		return
	}
	s.isShutdown = true
	s.mu.Unlock()

	log.Println("INFO: Iniciando shutdown de ImageUploadWorkerService...")

	// 1. Cerrar shutdown para señalizar a los workers que deben terminar
	close(s.shutdown)

	// 2. Cerrar jobQueue para que workers terminen cuando vacíen la cola
	// Los workers detectarán esto en su select y terminarán ordenadamente
	close(s.jobQueue)

	// 3. Esperar a que todos los workers terminen procesando trabajos pendientes
	log.Println("INFO: Esperando a que workers terminen...")
	s.wg.Wait()

	// 4. Cerrar resultQueue después de que todos los workers terminaron
	close(s.resultQueue)

	log.Println("INFO: ImageUploadWorkerService shutdown completado")
}
