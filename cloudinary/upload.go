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

type UploadOptions struct {
	Timestamp *string `json:"timestamp,omitempty"`
	//AccessControl           interface{} `json:"access_control,omitempty"`
	AccessMode     *string  `json:"access_mode,omitempty"`
	AllowedFormats *string  `json:"allowed_formats,omitempty"`
	Async          *bool    `json:"async,omitempty"`
	AutoTagging    *float64 `json:"auto_tagging,omitempty"`

	BackgroundRemoval *string `json:"background_removal,omitempty"`
	Backup            *bool   `json:"backup,omitempty"`

	Callback          *string `json:"callback,omitempty"`
	Categorization    *string `json:"categorization,omitempty"`
	Colors            *bool   `json:"colors,omitempty"`
	Context           *string `json:"context,omitempty"`
	CustomCoordinates *string `json:"custom_coordinates,omitempty"`

	Detection               *string `json:"detection,omitempty"`
	DiscardOriginalFilename *bool   `json:"discard_original_filename,omitempty"`

	Eager                *string `json:"eager,omitempty"`
	EagerAsync           *bool   `json:"eager_async,omitempty"`
	EagerNotificationURL *string `json:"eager_notification_url,omitempty"`
	Exif                 *bool   `json:"exif,omitempty"`

	FaceCoordinates *string `json:"face_coordinates,omitempty"`
	Faces           *bool   `json:"faces,omitempty"`
	Folder          *string `json:"folder,omitempty"`
	Format          *string `json:"format,omitempty"`

	Headers *string `json:"headers,omitempty"`

	ImageMetadata *bool `json:"image_metadata,omitempty"`
	Invalidate    *bool `json:"invalidate,omitempty"`

	Moderation *string `json:"moderation,omitempty"`

	NotificationURL *string `json:"notification_url,omitempty"`

	OCR       *string `json:"ocr,omitempty"`
	Overwrite *bool   `json:"overwrite,omitempty"`

	Phash    *bool   `json:"phash,omitempty"`
	Proxy    *string `json:"proxy,omitempty"`
	PublicId *string `json:"public_id,omitempty"`

	QualityAnalysis *bool `json:"quality_analysis,omitempty"`

	RawConvert        *string `json:"raw_convert,omitempty"`
	ReturnDeleteToken *bool   `json:"return_delete_token,omitempty"`

	Tags           *string `json:"tags,omitempty"`
	Transformation *string `json:"transformation,omitempty"`
	Type           *string `json:"type,omitempty"`

	UniqueFilename *bool `json:"unique_filename,omitempty"`
	// Upload Preset is required for unsigned uploading and
	// optional for signed uploading
	UploadPreset *string `json:"upload_preset,omitempty"`
	UseFilename  *bool   `json:"use_filename,omitempty"`

	ResourceType *string `json:"resource_type,omitempty"`
	//ResponsiveBreakpoints interface{} `json:"responsive_breakpoints,omitempty"`

	isUnsignedUpload bool `json:"-"`
}

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

type Opt func(uo *UploadOptions)

func WithUploadPreset(uploadPreset string) Opt {
	return func(uo *UploadOptions) {
		uo.UploadPreset = &uploadPreset
	}
}

func WithPublicId(id string) Opt {
	return func(uo *UploadOptions) {
		uo.PublicId = &id
	}
}

func WithFolder(folder string) Opt {
	return func(uo *UploadOptions) {
		uo.Folder = &folder
	}
}

func WithUseFilename(isUseFilename bool) Opt {
	return func(uo *UploadOptions) {
		uo.UseFilename = &isUseFilename
	}
}

func WithUniqueFilename(isUniqueFilename bool) Opt {
	return func(uo *UploadOptions) {
		uo.UniqueFilename = &isUniqueFilename
	}
}

func WithResourceType(resourceType string) Opt {
	return func(uo *UploadOptions) {
		uo.ResourceType = &resourceType
	}
}

func WithType(typeStr string) Opt {
	return func(uo *UploadOptions) {
		uo.Type = &typeStr
	}
}

func WithAccessMode(accessMode string) Opt {
	return func(uo *UploadOptions) {
		uo.AccessMode = &accessMode
	}
}

func WithDiscardOriginalFilename(dof bool) Opt {
	return func(uo *UploadOptions) {
		uo.DiscardOriginalFilename = &dof
	}
}

func WithOverwrite(isOverwrite bool) Opt {
	return func(uo *UploadOptions) {
		uo.Overwrite = &isOverwrite
	}
}

func WithTags(tags string) Opt {
	return func(uo *UploadOptions) {
		uo.Tags = &tags
	}
}

func WithContext(ctx string) Opt {
	return func(uo *UploadOptions) {
		uo.Context = &ctx
	}
}

func WithColors(hasColor bool) Opt {
	return func(uo *UploadOptions) {
		uo.Colors = &hasColor
	}
}

func WithFaces(returnFaces bool) Opt {
	return func(uo *UploadOptions) {
		uo.Faces = &returnFaces
	}
}

func WithQualityAnalysis(returnQualityAnalysis bool) Opt {
	return func(uo *UploadOptions) {
		uo.QualityAnalysis = &returnQualityAnalysis
	}
}

func WithImageMetadata(returnImageMetadata bool) Opt {
	return func(uo *UploadOptions) {
		uo.ImageMetadata = &returnImageMetadata
	}
}

func WithPhash(returnPhash bool) Opt {
	return func(uo *UploadOptions) {
		uo.Phash = &returnPhash
	}
}

func WithAutoTagging(autoTagging float64) Opt {
	return func(uo *UploadOptions) {
		uo.AutoTagging = &autoTagging
	}
}

func WithCategorization(c string) Opt {
	return func(uo *UploadOptions) {
		uo.Categorization = &c
	}
}

func WithDetection(d string) Opt {
	return func(uo *UploadOptions) {
		uo.Detection = &d
	}
}

func WithOCR(ocr string) Opt {
	return func(uo *UploadOptions) {
		uo.OCR = &ocr
	}
}

func WithExif(e bool) Opt {
	return func(uo *UploadOptions) {
		uo.Exif = &e
	}
}

func (uo *UploadOptions) GetPublicId() string {
	if uo.PublicId != nil {
		return *uo.PublicId
	}
	return ""
}

func (uo *UploadOptions) GetUploadPreset() string {
	if uo.UploadPreset != nil {
		return *uo.UploadPreset
	}
	return ""
}

func (us *UploadService) UploadImage(ctx context.Context, filePath string, opts ...Opt) (ur *UploadResponse, r *Response, err error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, nil, errors.New("invalid file")
	}
	opt := new(UploadOptions)
	for _, o := range opts {
		o(opt)
	}
	opt.isUnsignedUpload = false

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
		//return us.uploadFromURL(ctx, u, request, opt)
	}

	return ur, r, err
}

func (us *UploadService) UnsignedUploadImage(ctx context.Context, filePath string, uploadPreset string, opts ...Opt) (ur *UploadResponse, r *Response, err error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, nil, errors.New("invalid file")
	}
	if strings.TrimSpace(uploadPreset) == "" {
		return nil, nil, errors.New("upload_preset is required for unsigned uploading")
	}

	opt := new(UploadOptions)
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
		//return us.uploadFromURL(ctx, u, request, opt)
	}

	return &UploadResponse{}, &Response{}, nil
}

func (us *UploadService) uploadFromURL(ctx context.Context, url string, request *UploadRequest, opt *UploadOptions) (*UploadResponse, *Response, error) {
	req, err := us.client.NewRequest("POST", url, request)
	if err != nil {
		return nil, nil, err
	}

	ur := new(UploadResponse)
	resp, err := us.client.Do(ctx, req, ur)
	if err != nil {
		return nil, resp, err
	}

	return ur, resp, nil
}

func (us *UploadService) handleUploadFromLocalPath(ctx context.Context, u string, filePath string, opts *UploadOptions) (ur *UploadResponse, resp *Response, err error) {
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

func (us *UploadService) uploadFromS3(ctx context.Context, url string, request *UploadRequest, opt *UploadOptions) (*UploadResponse, *Response, error) {
	return &UploadResponse{}, &Response{}, nil
}

func (us *UploadService) uploadFromGoogleStorage(ctx context.Context, url string, request *UploadRequest, opt *UploadOptions) (*UploadResponse, *Response, error) {
	return &UploadResponse{}, &Response{}, nil
}

func (us *UploadService) buildParamsFromOptions(opts *UploadOptions, writer *multipart.Writer) error {
	if !opts.isUnsignedUpload {
		// Write timestamp
		timestamp := *opts.Timestamp
		ts, err := writer.CreateFormField("timestamp")
		if err != nil {
			return err
		}
		ts.Write([]byte(timestamp))

		// Write signature
		hash := sha1.New()
		part := fmt.Sprintf("timestamp=%s%s", timestamp, us.client.apiSecret)
		io.WriteString(hash, part)
		signature := fmt.Sprintf("%x", hash.Sum(nil))

		si, err := writer.CreateFormField("signature")
		if err != nil {
			return err
		}
		si.Write([]byte(signature))
	}
	var optMap map[string]interface{}
	optByte, _ := json.Marshal(opts)
	err := json.Unmarshal(optByte, &optMap)
	if err != nil {
		return err
	}

	for field, val := range optMap {
		valStr := fmt.Sprintf("%v", val)
		err := writer.WriteField(field, valStr)
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
