package access

import "github.com/Muvi7z/chat-auth-s/internal/model"

type Service struct {
	accessibleRoles map[string]int32
}

func NewService() *Service {
	accessibleRoles := make(map[string]int32)
	// Лезем в базу за данными о доступных ролях для каждого эндпоинта
	// Можно кэшировать данные, чтобы не лезть в базу каждый раз

	// Например, для эндпоинта /note_v1.NoteV1/Get доступна только роль admin
	accessibleRoles[model.ExamplePath] = 1

	return &Service{
		accessibleRoles: accessibleRoles,
	}
}
