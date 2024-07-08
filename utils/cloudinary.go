package utils

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// credentials initializes and returns a Cloudinary client and context
func credentials() (*cloudinary.Cloudinary, context.Context, error) {
	name := os.Getenv("CLOUDINARY_NAME")
	api_key := os.Getenv("CLOUDINARY_API")
	api_secret := os.Getenv("CLOUDINARY_SECRET")
	cld, err := cloudinary.NewFromParams(name, api_key, api_secret)
	if err != nil {
		return nil, nil, err
	}
	cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx, nil
}

func Uploader(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	cld, ctx, error := credentials()

	if error != nil {
		panic("Cloudinary Config Issues")
	}

	resp, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID:       "gharbheti/" + file.Filename,
		UniqueFilename: api.Bool(true),
		Overwrite:      api.Bool(false),
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

