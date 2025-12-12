package util

import "github.com/simesaba80/toybox-back/internal/infrastructure/config"

// ExtractS3KeyFromURL extracts S3 key from URL
// URL format: config.S3_BASE_URL + "/" + config.S3_BUCKET + "/" + key
// Example: "https://endpoint/bucket/dir/image/uuid/origin.png" -> "dir/image/uuid/origin.png"
func ExtractS3KeyFromURL(url string) string {
	prefix := config.S3_BASE_URL + "/" + config.S3_BUCKET + "/"
	if len(url) <= len(prefix) {
		return ""
	}
	if url[:len(prefix)] != prefix {
		return ""
	}
	return url[len(prefix):]
}
