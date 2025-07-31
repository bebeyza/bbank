package models

// All models in one place for easy migration
func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&Balance{},
		&Transaction{},
		&AuditLog{},
	}
}
