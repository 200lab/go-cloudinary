package cloudinary

import "encoding/base64"

func getBase64EncodedString(u string, p string) string {
	return base64.StdEncoding.EncodeToString([]byte(u + ":" + p))
}
