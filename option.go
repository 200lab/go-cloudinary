package cloudinary

type ResourceType string

type Options struct {
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

	KeepOriginal *bool `json:"keep_original,omitempty"`

	Moderation *string `json:"moderation,omitempty"`

	NextCursor      *string `json:"next_cursor,omitempty"`
	NotificationURL *string `json:"notification_url,omitempty"`

	OCR       *string `json:"ocr,omitempty"`
	Overwrite *bool   `json:"overwrite,omitempty"`

	Phash    *bool   `json:"phash,omitempty"`
	Proxy    *string `json:"proxy,omitempty"`
	PublicId *string `json:"public_id,omitempty"`

	QualityAnalysis *bool `json:"quality_analysis,omitempty"`

	RawConvert        *string `json:"raw_convert,omitempty"`
	ReturnDeleteToken *bool   `json:"return_delete_token,omitempty"`
	ResourceType      *string `json:"resource_type,omitempty"`
	//ResponsiveBreakpoints interface{} `json:"responsive_breakpoints,omitempty"`

	Tags            *string `json:"tags,omitempty"`
	Timestamp       *string `json:"timestamp,omitempty"`
	Transformations *string `json:"transformation,omitempty"`
	Type            *string `json:"type,omitempty"`

	UniqueFilename *bool `json:"unique_filename,omitempty"`
	// Upload Preset is required for unsigned uploading and
	// optional for signed uploading
	UploadPreset *string `json:"upload_preset,omitempty"`
	UseFilename  *bool   `json:"use_filename,omitempty"`

	isUnsignedUpload bool
}

type SetOpts func(opts *Options)

func WithUploadPreset(uploadPreset string) SetOpts {
	return func(uo *Options) {
		uo.UploadPreset = &uploadPreset
	}
}

func WithPublicId(id string) SetOpts {
	return func(uo *Options) {
		uo.PublicId = &id
	}
}

func WithFolder(folder string) SetOpts {
	return func(uo *Options) {
		uo.Folder = &folder
	}
}

func WithUseFilename(isUseFilename bool) SetOpts {
	return func(uo *Options) {
		uo.UseFilename = &isUseFilename
	}
}

func WithUniqueFilename(isUniqueFilename bool) SetOpts {
	return func(uo *Options) {
		uo.UniqueFilename = &isUniqueFilename
	}
}

func WithType(typeStr string) SetOpts {
	return func(uo *Options) {
		uo.Type = &typeStr
	}
}

func WithAccessMode(accessMode string) SetOpts {
	return func(uo *Options) {
		uo.AccessMode = &accessMode
	}
}

func WithDiscardOriginalFilename(dof bool) SetOpts {
	return func(uo *Options) {
		uo.DiscardOriginalFilename = &dof
	}
}

func WithOverwrite(isOverwrite bool) SetOpts {
	return func(uo *Options) {
		uo.Overwrite = &isOverwrite
	}
}

func WithTags(tags string) SetOpts {
	return func(uo *Options) {
		uo.Tags = &tags
	}
}

func WithContext(ctx string) SetOpts {
	return func(uo *Options) {
		uo.Context = &ctx
	}
}

func WithColors(hasColor bool) SetOpts {
	return func(uo *Options) {
		uo.Colors = &hasColor
	}
}

func WithFaces(returnFaces bool) SetOpts {
	return func(uo *Options) {
		uo.Faces = &returnFaces
	}
}

func WithQualityAnalysis(returnQualityAnalysis bool) SetOpts {
	return func(uo *Options) {
		uo.QualityAnalysis = &returnQualityAnalysis
	}
}

func WithImageMetadata(returnImageMetadata bool) SetOpts {
	return func(uo *Options) {
		uo.ImageMetadata = &returnImageMetadata
	}
}

func WithPhash(returnPhash bool) SetOpts {
	return func(uo *Options) {
		uo.Phash = &returnPhash
	}
}

func WithAutoTagging(autoTagging float64) SetOpts {
	return func(uo *Options) {
		uo.AutoTagging = &autoTagging
	}
}

func WithCategorization(c string) SetOpts {
	return func(uo *Options) {
		uo.Categorization = &c
	}
}

func WithDetection(d string) SetOpts {
	return func(uo *Options) {
		uo.Detection = &d
	}
}

func WithOCR(ocr string) SetOpts {
	return func(uo *Options) {
		uo.OCR = &ocr
	}
}

func WithExif(e bool) SetOpts {
	return func(uo *Options) {
		uo.Exif = &e
	}
}

func (uo *Options) GetPublicId() string {
	if uo.PublicId != nil {
		return *uo.PublicId
	}
	return ""
}

func (uo *Options) GetUploadPreset() string {
	if uo.UploadPreset != nil {
		return *uo.UploadPreset
	}
	return ""
}

func (uo *Options) GetTimestamp() string {
	if uo.Timestamp != nil {
		return *uo.Timestamp
	}
	return ""
}

func (uo *Options) GetResourceType() string {
	if uo.ResourceType != nil {
		return *uo.ResourceType
	}
	return ""
}

func (uo *Options) GetType() string {
	if uo.Type != nil {
		return *uo.Type
	}
	return ""
}
func WithResourceType(resourceType string) SetOpts {
	return func(opts *Options) {
		opts.ResourceType = &resourceType
	}
}
