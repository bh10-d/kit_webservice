package util

import (
	"strings"
	"math"
	// "archive/tar"
	// "compress/gzip"
	// "fmt"
	// "bytes"
	"go-controller/db"
	"go-controller/dto"
	"go-controller/model"
	"os"
	"gorm.io/gorm"
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

// GetPaginationDefaults returns default pagination values
func GetPaginationDefaults(req dto.PaginationRequest) dto.PaginationRequest {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}
	if req.Sort == "" {
		req.Sort = "id"
	}
	if req.Order == "" {
		req.Order = "desc"
	}
	return req
}

// CalculatePaginationMeta calculates pagination metadata
func CalculatePaginationMeta(page, pageSize int, total int64) dto.PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	
	return dto.PaginationMeta{
		Page:        page,
		PageSize:    pageSize,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}

// ApplyPagination applies pagination, sorting and searching to a GORM query
func ApplyPagination(query *gorm.DB, req dto.PaginationRequest) *gorm.DB {
	req = GetPaginationDefaults(req)
	
	// Apply sorting
	orderClause := req.Sort + " " + req.Order
	query = query.Order(orderClause)
	
	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)
	
	return query
}

// ApplySearch applies search filter to queries based on model type
func ApplySearch(query *gorm.DB, search string, modelType string) *gorm.DB {
	if search == "" {
		return query
	}
	
	searchPattern := "%" + search + "%"
	
	switch modelType {
	case "job":
		return query.Where("runner_id ILIKE ? OR msg_id ILIKE ? OR status ILIKE ?", 
			searchPattern, searchPattern, searchPattern)
	case "runner":
		return query.Where("id ILIKE ? OR host_name ILIKE ? OR ip ILIKE ? OR tags ILIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern)
	case "script":
		return query.Where("script_id ILIKE ? OR file_name ILIKE ? OR description ILIKE ?", 
			searchPattern, searchPattern, searchPattern)
	case "log":
		return query.Where("msg_id ILIKE ? OR runner_id ILIKE ? OR status ILIKE ? OR message ILIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern)
	default:
		return query
	}
}