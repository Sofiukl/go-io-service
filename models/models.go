package models

// GenericResponse - Common Response structure
type GenericResponse struct {
	Error   bool        `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Details string      `json:"details,omitempty"`
}

// V2ApplicationMigration model
type V2ApplicationMigration struct {
	ID    string `json:"_id" bson:"_id"`
	State string `json:"state,omitempty" bson:"state"`
}

//UploadHistory - model for upload history
type UploadHistory struct {
	FileName string `bson:"file_name,omitempty"`
	FileKey  string `bson:"file_key,omitempty"`
}
