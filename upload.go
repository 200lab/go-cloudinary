package cloudinary

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// UploadService handles communication with the uploading related

type UploadService service

type UploadRequest struct {
	// Required fields to call upload request
	File string `json:"file"`
	//Timestamp string `json:"timestamp"`
	//UploadPreset *string `json:"upload_preset, omitempty"`
}

//type Options struct {
//	isUnsignedUpload bool `json:"-"`
//}

type UploadResponse struct {
	PublicId         string   `json:"public_id"`
	Version          int64    `json:"version"`
	Signature        string   `json:"signature"`
	Width            int64    `json:"width"`
	Height           int64    `json:"height"`
	Format           string   `json:"format"`
	ResourceType     string   `json:"resource_type"`
	CreatedAt        string   `json:"created_at"`
	Tags             []string `json:"tags"`
	Bytes            int64    `json:"bytes"`
	Type             string   `json:"type"`
	Etag             string   `json:"etag"`
	Placeholder      bool     `json:"placeholder"`
	URL              string   `json:"url"`
	SecureURL        string   `json:"secure_url"`
	AccessMode       string   `json:"access_mode"`
	OriginalFilename string   `json:"original_filename"`
}

// UploadImage handle signed uploading image to Cloudinary
// Signed request are required `signature` parameters
func (us *UploadService) UploadImage(ctx context.Context, filePath string, opts ...SetOpts) (ur *UploadResponse, r *Response, err error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, nil, errors.New("invalid file")
	}
	opt := new(Options)
	for _, o := range opts {
		o(opt)
	}
	opt.isUnsignedUpload = false

	u := fmt.Sprintf("image/upload")

	switch {
	case strings.Contains(filePath, "://"):
		// Upload image using HTTPS URL or HTTP
		return us.uploadFromURL(ctx, u, filePath, opt)
	case strings.Contains(filePath, "s3"):
		// Upload image using Amazon S3
		//return us.uploadFromS3(ctx, u, request, opt)
	case strings.Contains(filePath, "gs"):
		// Upload image using Google Storage
		//return us.uploadFromGoogleStorage(ctx, u, request, opt)

	default:
		return us.handleUploadFromLocalPath(ctx, u, filePath, opt)
	}

	return ur, r, err
}

// UnsignedUploadImage handle unsigned uploading image to Cloudinary.
// Unsigned request are restricted to the following allowed parameters:
// `public_id`, `folder`, `callback`, `tags`, `context`,
// `face_coordinates` (image-only), `custom_coordinates` (image-only) and `upload_preset`.
// Most of the other upload parameters can be defined in your `upload_preset`.
//
// Additionally, although the `public_id` parameter can be specified,
// the `overwrite` parameter is always set to `false` for unsigned uploads
// to prevent overwriting existing file.
func (us *UploadService) UnsignedUploadImage(ctx context.Context, filePath string, uploadPreset string, opts ...SetOpts) (ur *UploadResponse, r *Response, err error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, nil, errors.New("invalid file")
	}
	if strings.TrimSpace(uploadPreset) == "" {
		return nil, nil, errors.New("upload_preset is required for unsigned uploading")
	}

	opt := new(Options)
	for _, o := range opts {
		o(opt)
	}
	opt.isUnsignedUpload = true
	opt.UploadPreset = &uploadPreset

	u := fmt.Sprintf("image/upload")

	switch {
	case strings.HasPrefix(filePath, "/"):
		// Upload image using local path
		return us.handleUploadFromLocalPath(ctx, u, filePath, opt)
	case strings.HasPrefix(filePath, "s3"):
		// Upload image using Amazon S3
		//return us.uploadFromS3(ctx, u, request, opt)
	case strings.HasPrefix(filePath, "gs"):
		// Upload image using Google Storage
		//return us.uploadFromGoogleStorage(ctx, u, request, opt)

	default:
		// Upload image using HTTPS URL or HTTP
		return us.uploadFromURL(ctx, u, filePath, opt)
	}

	return &UploadResponse{}, &Response{}, nil
}

func (us *UploadService) uploadFromURL(ctx context.Context, u, fileURL string, opts *Options) (ur *UploadResponse, resp *Response, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("file", fileURL); err != nil {
		return ur, resp, err
	}

	if !opts.isUnsignedUpload {
		timestamp := strconv.Itoa(int(time.Now().UTC().Unix()))
		opts.Timestamp = &timestamp

		ak, err := writer.CreateFormField("api_key")
		if err != nil {
			return ur, resp, err
		}
		_, err = ak.Write([]byte(us.client.apiKey))
		if err != nil {
			return ur, resp, err
		}
	}

	if opts != nil {
		if err := us.buildParamsFromOptions(opts, writer); err != nil {
			return ur, resp, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, nil, err
	}

	req, err := us.client.NewUploadRequest(u, body, writer)
	if err != nil {
		return nil, nil, err
	}

	ur = new(UploadResponse)
	resp, err = us.client.Do(ctx, req, ur)
	if err != nil {
		return nil, resp, err
	}

	return ur, resp, nil
}

func (us *UploadService) handleUploadFromLocalPath(ctx context.Context, u, filePath string, opts *Options) (ur *UploadResponse, resp *Response, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, _, err := us.openFile(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	if stat.IsDir() {
		return nil, nil, errors.New("the asset to upload can't be a directory")
	}

	if !opts.isUnsignedUpload {
		timestamp := strconv.Itoa(int(time.Now().UTC().Unix()))
		opts.Timestamp = &timestamp

		ak, err := writer.CreateFormField("api_key")
		if err != nil {
			return ur, resp, err
		}
		_, err = ak.Write([]byte(us.client.apiKey))
		if err != nil {
			return ur, resp, err
		}
	}

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, nil, err
	}
	_, err = io.Copy(part, file)

	if opts != nil {
		if err := us.buildParamsFromOptions(opts, writer); err != nil {
			return nil, nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, nil, err
	}

	req, err := us.client.NewUploadRequest(u, body, writer)
	if err != nil {
		return nil, nil, err
	}

	ur = new(UploadResponse)
	resp, err = us.client.Do(ctx, req, ur)
	if err != nil {
		return nil, resp, err
	}

	return ur, resp, nil
}

func (us *UploadService) uploadFromS3(ctx context.Context, url string, request *UploadRequest, opt *Options) (*UploadResponse, *Response, error) {
	return &UploadResponse{}, &Response{}, nil
}

func (us *UploadService) uploadFromGoogleStorage(ctx context.Context, url string, request *UploadRequest, opt *Options) (*UploadResponse, *Response, error) {
	return &UploadResponse{}, &Response{}, nil
}

func (us *UploadService) buildParamsFromOptions(opts *Options, writer *multipart.Writer) error {
	if !opts.isUnsignedUpload {
		// Write timestamp
		timestamp := opts.GetTimestamp()
		ts, err := writer.CreateFormField("timestamp")
		if err != nil {
			return err
		}
		_, err = ts.Write([]byte(timestamp))
		if err != nil {
			return err
		}
	}
	var optMap map[string]interface{}
	optByte, _ := json.Marshal(opts)
	err := json.Unmarshal(optByte, &optMap)
	if err != nil {
		return err
	}

	hash := sha1.New()
	params := make([]string, 0)

	for field, val := range optMap {
		valStr := fmt.Sprintf("%v", val)
		err := writer.WriteField(field, valStr)
		if err != nil {
			return err
		}

		params = append(params, fmt.Sprintf("%s=%s", field, valStr))
	}

	if !opts.isUnsignedUpload {
		part := strings.Join(params, "&")
		hash.Write([]byte(part + us.client.apiSecret))
		signature := fmt.Sprintf("%x", hash.Sum(nil))

		si, err := writer.CreateFormField("signature")
		if err != nil {
			return err
		}
		_, err = si.Write([]byte(signature))
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *UploadService) openFile(filePath string) (file *os.File, dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return nil, dir, err
	}
	file, err = os.Open(dir + filePath)
	return file, dir, err

}

func (us *UploadService) getFilename(filePath string) string {
	var extension = filepath.Ext(filePath)
	return filePath[0 : len(filePath)-len(extension)]
}
