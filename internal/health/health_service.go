package health

import "audiscript_be/database"

type Service interface {
    // Check trả về map đã format sẵn
    Check() map[string]interface{}
}

type service struct {
    db database.Service
}

func NewService(db database.Service) Service {
    return &service{db: db}
}

func (s *service) Check() map[string]interface{} {
    result := make(map[string]interface{})
    for k, v := range s.db.Health() {
        result[k] = v
    }
    return result
}