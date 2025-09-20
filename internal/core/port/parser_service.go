package port

import "cleo.com/internal/core/domain/model"

//go:generate moq -pkg mocks -out ./mocks/parser_service.go . HealthMetricParserService

type HealthMetricParserService interface {
	ParseClinicalNote(note *model.ClinicalNote) (*model.HealthMetric, error)
}
