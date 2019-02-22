package cloudinary

type AdminService service

type AdminOptions struct {
}

type AdminResponse struct{}

type AdminOpt func(ao *AdminOptions)

func (as *AdminService) DeleteResources(publicId []string, opts ...AdminOpt) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, err
}

func (as *AdminService) DeleteResourcesByPrefix(prefix string, opts ...AdminOpt) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}

func (as *AdminService) DeleteAllResources(opts ...AdminOpt) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}

func (as *AdminService) DeleteResourcesByTag(tag string, opts ...AdminOpt) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}

func (as *AdminService) DeleteDerivedResources(derivedResourceIds, opts ...AdminOpt) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}

func (as *AdminService) DeleteDerivedResourcesByTransformation(publicId, transformation []string, opts ...AdminOpt) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}
