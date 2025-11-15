package services

type ICloudinaryService interface {
	UploadImage(localPath string) (string, error)
}
