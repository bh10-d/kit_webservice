package service

import (
	// "encoding/json"
	"errors"
	"fmt"
	// "github.com/google/uuid"
	"go-controller/internal/database"
	"go-controller/internal/dto"
	"go-controller/internal/model"
	// "fmt"
)

// ScriptService handles script-related operations
type ScriptService struct{}

// NewScriptService creates a new script service
func NewScriptService() *ScriptService {
	return &ScriptService{}
}

// GetAllScripts returns all available scripts
func (s *ScriptService) GetAllScripts() ([]model.Scripts, error) {
	var scripts []model.Scripts
	if err := db.DB.Find(&scripts).Error; err != nil {
		return nil, err
	}
	return scripts, nil
}

// // GetScriptByID returns a script by its ID
func (s *ScriptService) GetScriptByID(scriptID string) (*model.Scripts, error) {
	var script model.Scripts
	if err := db.DB.Where("script_id = ?", scriptID).First(&script).Error; err != nil {
		return nil, err
	}
	return &script, nil
}

// GetScriptsByOperation returns scripts for a specific operation type
// func (s *ScriptService) GetScriptsByOperation(operation string) ([]model.Scripts, error) {
// 	var scripts []model.Scripts
// 	if err := db.DB.Where("description LIKE ?", "%"+operation+"%").Find(&scripts).Error; err != nil {
// 		return nil, err
// 	}
// 	return scripts, nil
// }

// BuildScriptPayload builds the payload for script execution
func (s *ScriptService) BuildScriptPayload(scriptID string, parameters map[string]interface{}) (map[string]interface{}, error) {
	script, err := s.GetScriptByID(scriptID)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"script": script.FileName,
	}

	tmpParameters := map[string]interface{}{}

	// Add parameters to payload
	for key, value := range parameters {
		tmpParameters[key] = value
	}

	

	// Validate required parameters
	if err := s.validateParameters(script, tmpParameters); err != nil {
		return nil, err
	}

	payload["parameters"] = tmpParameters

	// fmt.Println("Built payload:", payload)

	return payload, nil
}

// validateParameters validates if all required parameters are provided
func (s *ScriptService) validateParameters(script *model.Scripts, parameters map[string]interface{}) error {
	// Convert script.Param to required parameters list
	for _, paramName := range script.Param {
		if paramName != "" {
			if _, exists := parameters[paramName]; !exists {
				return errors.New("missing required parameter: " + paramName)
			}
		}
	}
	return nil
}

// GetScriptParameters returns parameter information for a script
func (s *ScriptService) GetScriptParameters(scriptID string) ([]dto.ScriptParameterInfo, error) {
	script, err := s.GetScriptByID(scriptID)
	if err != nil {
		return nil, err
	}

	var params []dto.ScriptParameterInfo
	for _, paramName := range script.Param {
		if paramName != "" {
			param := dto.ScriptParameterInfo{
				Name:        paramName,
				Type:        "string", // Default type, could be enhanced
				Required:    true,     // Default to required, could be enhanced
				Description: "Parameter for " + paramName,
			}
			params = append(params, param)
		}
	}

	return params, nil
}

// InitializeDefaultScripts creates default scripts in database
// func (s *ScriptService) InitializeDefaultScripts() error {
// 	defaultScripts := []model.Scripts{
// 		{
// 			ScriptID:    "check-site",
// 			FileName:    "check_site.sh",
// 			Description: "Check if site exists",
// 			Param:       []string{"subDomain"},
// 		},
// 		{
// 			ScriptID:    "create-site",
// 			FileName:    "create_site.sh", 
// 			Description: "Create new site",
// 			Param:       []string{"subDomain"},
// 		},
// 		{
// 			ScriptID:    "remove-site",
// 			FileName:    "remove_site.sh",
// 			Description: "Remove existing site", 
// 			Param:       []string{"subDomain"},
// 		},
// 		{
// 			ScriptID:    "update-site",
// 			FileName:    "update_site.sh",
// 			Description: "Update site configuration",
// 			Param:       []string{"oldSubDomain", "newSubDomain"},
// 		},
// 	}

// 	for _, script := range defaultScripts {
// 		// Check if script already exists
// 		var existingScript model.Scripts
// 		if err := db.DB.Where("script_id = ?", script.ScriptID).First(&existingScript).Error; err != nil {
// 			// Script doesn't exist, create it
// 			if err := db.DB.Create(&script).Error; err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// ConvertLegacyRequest converts old-style requests to script execution requests
// func (s *ScriptService) ConvertSiteRequestToScriptRequest(req dto.SiteRequest, operation string) (*dto.ScriptExecutionRequest, error) {
// 	var scriptID string
	
// 	switch operation {
// 	case "check":
// 		scriptID = "check-site"
// 	case "create":
// 		scriptID = "create-site" 
// 	case "remove":
// 		scriptID = "remove-site"
// 	default:
// 		return nil, errors.New("unsupported operation: " + operation)
// 	}

// 	parameters := map[string]interface{}{
// 		"subDomain": req.SubDomain,
// 	}

// 	return &dto.ScriptExecutionRequest{
// 		ScriptID:   scriptID,
// 		Parameters: parameters,
// 		Tag:        req.Tag,
// 	}, nil
// }


func (s *ScriptService) ConvertSiteRequestToScriptRequest(req dto.SiteRequest, operation string) (*dto.ScriptExecutionRequest, error) {
	parameters := map[string]interface{}{
		"subDomain": req.SubDomain,
	}

	return &dto.ScriptExecutionRequest{
		ScriptID:   req.ScriptID, // ✅ giữ nguyên int từ DTO
		// ScriptID:   uuid.New().String(), // ✅ giữ nguyên int từ DTO
		Parameters: parameters,
		Tag:        req.Tag,
		// Operation:  operation,   // gắn thêm nếu muốn phân biệt
	}, nil
}


func (s *ScriptService) CreateScript(script *model.Scripts) error {
	// Đảm bảo arrays được khởi tạo đúng cách
	if script.Param == nil {
		script.Param = make([]string, 0)
	}
	if script.Tag == nil {
		script.Tag = make([]string, 0)
	}
	if script.Runner == nil {
		script.Runner = make([]string, 0)
	}
	
	return db.DB.Create(script).Error
}


func (s *ScriptService) UpdateScript(script *model.Scripts) error {
	fmt.Printf("UpdateScript - ID: %s\n", script.ScriptID)
	fmt.Printf("UpdateScript - Param: %v\n", script.Param)
	fmt.Printf("UpdateScript - Tag: %v\n", script.Tag)
	fmt.Printf("UpdateScript - Runner: %v\n", script.Runner)
	
	// Sử dụng Updates() thay vì Save() để xử lý PostgreSQL arrays đúng cách
	result := db.DB.Model(&model.Scripts{}).Where("script_id = ?", script.ScriptID).Updates(map[string]interface{}{
		"file_name":   script.FileName,
		"description": script.Description,
		"param":       script.Param,
		"status":      script.Status,
		"tag":         script.Tag,
		"runner":      script.Runner,
	})
	
	fmt.Printf("UpdateScript - Affected rows: %d\n", result.RowsAffected)
	if result.Error != nil {
		fmt.Printf("UpdateScript - Error: %v\n", result.Error)
	}
	
	return result.Error
}

func (s *ScriptService) UpdateScriptStatus(script *model.Scripts) error {
	return db.DB.Model(&model.Scripts{}).Where("script_id = ?", script.ScriptID).Update("status", script.Status).Error
}

func (s *ScriptService) DeleteScript(scriptID string) error {
	return db.DB.Where("script_id = ?", scriptID).Delete(&model.Scripts{}).Error
}