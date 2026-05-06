package utils

import (
	"fmt"
	"os"
	"strings"
)

var validExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

var screenshotPatterns = []string{
	"screenshots",
	"screenshot",
	"captura_de_pantalla",
	"captura",
	"screen shot",
	"screen capture",
	" print",
	"_print",
}

var NoScreenshots bool
var IgnoreExtensions []string
var AdditionalExtensions []string

func IsScreenshot(filename string) bool {
	lower := strings.ToLower(filename)
	for _, pattern := range screenshotPatterns {
		if strings.Contains(lower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func IsValidImage(path string) bool {
	filename := GetFilename(path)
	ext := strings.ToLower(getExt(filename))

	isValid := validExtensions[ext]
	for _, extra := range AdditionalExtensions {
		if ext == "."+extra || ext == extra {
			isValid = true
			break
		}
	}

	if !isValid {
		return false
	}

	for _, ignored := range IgnoreExtensions {
		if ext == "."+ignored || ext == ignored {
			return false
		}
	}

	if NoScreenshots && IsScreenshot(filename) {
		return false
	}
	return true
}

func getExt(filename string) string {
	i := strings.LastIndex(filename, ".")
	if i == -1 {
		return ""
	}
	return filename[i:]
}

func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetFilename(path string) string {
	i := strings.LastIndex(path, "/")
	if i == -1 {
		return path
	}
	return path[i+1:]
}

func GetFilebase(filename string) string {
	i := strings.LastIndex(filename, ".")
	if i == -1 {
		return filename
	}
	return filename[:i]
}

func GetExtension(filename string) string {
	ext := getExt(filename)
	return strings.ToLower(ext)
}

func GetDestPath(destDir, filename string, mode string) (string, error) {
	base := GetFilebase(filename)
	ext := GetExtension(filename)
	destPath := destDir + "/" + filename

	if !FileExists(destPath) {
		return destPath, nil
	}

	if mode == "skip" {
		return "", nil
	}

	for i := 1; ; i++ {
		newFilename := fmt.Sprintf("%s_%d%s", base, i, ext)
		newPath := destDir + "/" + newFilename
		if !FileExists(newPath) {
			return newPath, nil
		}
	}
}