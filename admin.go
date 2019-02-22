package cloudinary

type AdminService service

type AdminOptions struct{}

type AdminResponse struct{}

type SetAdminOpts func(ao *AdminOptions)

// DeleteResource deletes all resources with the given publicIds
// publicIds is a array that store up to 100 ids
func (as *AdminService) DeleteResources(publicIds []string, opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
	o := new(Options)
	for _, setOptions := range opts {
		setOptions(o)
	}

	return &AdminResponse{}, &Response{}, err
}

func (as *AdminService) DeleteResourcesByPrefix(prefix string, opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}

func (as *AdminService) DeleteAllResources(opts ...SetAdminOpts) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}

func (as *AdminService) DeleteResourcesByTag(tag string, opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}

func (as *AdminService) DeleteDerivedResources(derivedResourceIds string, opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}

func (as *AdminService) DeleteDerivedResourcesByTransformation(publicId, transformation []string, opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}
