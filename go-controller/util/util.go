package util

import (
	"strings"
	// "archive/tar"
	// "compress/gzip"
	// "fmt"
	// "bytes"
	"go-controller/db"
	// "go-controller/dto"
	"go-controller/model"
	"os"
	// "fmt"
	"errors"
)

func SplitTags(tags string) []string {
	var res []string
	for _, t := range SplitAndTrim(tags, ",") {
		if t != "" {
			res = append(res, t)
		}
	}
	return res
}

func SplitAndTrim(s, sep string) []string {
	var out []string
	for _, part := range strings.Split(s, sep) {
		out = append(out, strings.TrimSpace(part))
	}
	return out
}


// func ListScripts() []string {
// 	var scripts []string
// 	files, err := os.ReadDir("./scripts")
// 	if err != nil {
// 		return scripts
// 	}
// 	for _, file := range files {
// 		if !file.IsDir() {
// 			scripts = append(scripts, file.Name())
// 		}
// 	}
// 	return scripts
// }



func ListScripts() [] string {
	var scripts []string
	files, err := os.ReadDir("./scripts")
	if err != nil {
		return scripts
	}
	for _, file := range files {
		if !file.IsDir(){
			scripts = append(scripts, file.Name())
		}
	}
	return scripts
}



// TarGzFolder nén một thư mục thành file .tar.gz bằng lệnh tar của Linux
// Chỉ trả về log kết quả thành công hay lỗi
// func TarGzFolder(srcDir, destTarGz string) error {
// 	var stderr bytes.Buffer

// 	cmd := exec.Command("tar", "-czf", destTarGz, "-C", srcDir, ".")
// 	cmd.Stderr = &stderr

// 	if err := cmd.Run(); err != nil {
// 		return fmt.Errorf("compress failed: %v, detail: %s", err, stderr.String())
// 	}

// 	return nil
// }


// func ZipFolderWithPassword(srcDir, destZip, password string) error {
// 	var stderr bytes.Buffer
// 	// -r: recursive, -P: password
// 	cmd := exec.Command("zip", "-p", password, destZip, ".")
// 	cmd.Dir = srcDir
// 	cmd.Stderr = &stderr

// 	if err := cmd.Run(); err != nil {
// 		return fmt.Errorf("zip failed: %v, detail: %s", err, stderr.String())
// 	}
// 	return nil
// }




func CheckStatus (scriptID string) error {
	// fmt.Printf("Checking status for script ID: %s\n", scriptID)
	var script model.Scripts
	if err := db.DB.First(&script, "script_id = ?", scriptID).Error; err != nil {
		// return errors.New("Service not found")
		return errors.New(scriptID)
	}
	
	if !script.Status {
		return errors.New("Service is not active")

	}
	return nil
}