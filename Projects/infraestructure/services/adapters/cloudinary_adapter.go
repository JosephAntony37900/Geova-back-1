//Geova-back-1/Projects/infraestructure/services/adapters/cloudinary_adapter.go
package adapters

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryAdapter struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryAdapter() (*CloudinaryAdapter, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("error configurando cloudinary: %v", err)
	}

	return &CloudinaryAdapter{cld: cld}, nil
}

func (c *CloudinaryAdapter) UploadImage(localPath string) (string, error) {
	ctx := context.Background()
	uploadResult, err := c.cld.Upload.Upload(ctx, localPath, uploader.UploadParams{})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}
