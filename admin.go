package cloudinary

import (
	"context"
	"fmt"
)

type AdminService service

type deleteResourcesOptions struct {
	KeepOriginal bool   `json:"keep_original"`
	Invalidate   bool   `json:"invalidate"`
	NextCursor   string `json:"next_cursor"`
}

func (dro *deleteResourcesOptions) SetKeepOriginal(ko bool) {
	dro.KeepOriginal = ko
}

func (dro *deleteResourcesOptions) SetInvalidate(i bool) {
	dro.Invalidate = i
}

func (dro *deleteResourcesOptions) SetNextCursor(nc string) {
	dro.NextCursor = nc
}

type AdminResponse struct {
	Deleted interface{} `json:"deleted"`
	Partial bool        `json:"partial"`
}

// DeleteResource deletes all resources with the given publicIds
// publicIds is a array that store up to 100 ids

// Documentation: https://cloudinary.com/documentation/admin_api#delete_all_or_selected_resources
func (as *AdminService) DeleteResources(ctx context.Context, publicIds []string, opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
	o := new(Options)
	for _, setOptions := range opts {
		setOptions(o)
	}

	resourceType := o.GetResourceType()
	if resourceType == "" {
		resourceType = "image"
	}
	storageType := o.GetType()
	if storageType == "" {
		storageType = "upload"
	}

	u := fmt.Sprintf("/resources/%s/%s", resourceType, storageType)

	request, err := as.client.NewRequest("DELETE", u, nil)
	for _, publicId := range publicIds {
		request.URL.Query().Add("public_ids[]", publicId)
		request.URL.Query().Set("public_ids[]", publicId)
	}
	if err != nil {
		return &AdminResponse{}, &Response{}, err
	}

	resp, err = as.client.Do(ctx, request, ar)
	return ar, resp, err
}

func (as *AdminService) DeleteResourcesByPrefix(prefix string, opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
	return &AdminResponse{}, &Response{}, nil
}

func (as *AdminService) DeleteAllResources(opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
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
func (as *AdminService) withBasicAuthentication(request *http.Request) {
	encodedStr := getBase64EncodedString(as.client.apiKey, as.client.apiSecret)
	request.Header.Set("Authorization", "Basic "+encodedStr)
}
