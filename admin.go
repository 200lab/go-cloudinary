package cloudinary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

func (ar *AdminResponse) ToJSON() string {
	b, _ := json.Marshal(ar)
	return string(b)
}

// DeleteResource deletes all resources with the given publicIds
// publicIds is a array that store up to 100 ids
//
// Documentation: https://cloudinary.com/documentation/admin_api#delete_all_or_selected_resources
func (as *AdminService) DeleteResources(ctx context.Context, publicIds []string, opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
	o := new(Options)
	params := make(map[string]string)
	for _, setOptions := range opts {
		setOptions(o)
	}
	keepOriginal := o.GetKeepOriginal()
	if keepOriginal {
		params["keep_original"] = strconv.FormatBool(keepOriginal)
	}
	for _, pId := range publicIds {
		params["public_ids[]"] = pId
	}

	resourceType := o.GetResourceType()
	if resourceType == "" {
		resourceType = "image"
	}
	storageType := o.GetType()
	if storageType == "" {
		storageType = "upload"
	}

	u := fmt.Sprintf("resources/%s/%s", resourceType, storageType)
	u = as.buildURLStrWithParams(u, params)

	request, err := as.client.NewRequest("DELETE", u, o)
	if err != nil {
		return &AdminResponse{}, &Response{}, err
	}
	as.withBasicAuthentication(request)

	ar = new(AdminResponse)
	resp, err = as.client.Do(ctx, request, ar)
	return ar, resp, err
}

func (as *AdminService) DeleteResourcesByPrefix(ctx context.Context, prefix string, opts ...SetOpts) (ar *AdminResponse, resp *Response, err error) {
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

	u := fmt.Sprintf("resources/%s/%s?prefix=%s", resourceType, storageType, prefix)

	request, err := as.client.NewRequest("DELETE", u, o)
	if err != nil {
		return &AdminResponse{}, &Response{}, err
	}
	as.withBasicAuthentication(request)

	ar = new(AdminResponse)
	resp, err = as.client.Do(ctx, request, ar)
	return ar, resp, err
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

// buildURLStrWithParams is return url string that contain query string with the given parameters
func (as *AdminService) buildURLStrWithParams(u string, params map[string]string) string {
	urlObject, _ := url.Parse(u)
	q := urlObject.Query()

	for key, val := range params {
		q.Add(key, val)
	}

	urlObject.RawQuery = q.Encode()
	return urlObject.String()
}

// Cloudinary Admin API use Basic Authentication over secure HTTP.
// API_KEY and API_SECRET are used for the authentication
//
// Read more about Basic Authentication at link below:
// https://swagger.io/docs/specification/authentication/basic-authentication/
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication
func (as *AdminService) withBasicAuthentication(request *http.Request) {
	encodedStr := getBase64EncodedString(as.client.apiKey, as.client.apiSecret)
	request.Header.Set("Authorization", "Basic "+encodedStr)
}
