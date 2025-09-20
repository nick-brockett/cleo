package service_test

import (
	"errors"
	"testing"

	"cleo.com/internal/core/domain/model"
	"cleo.com/internal/core/service"
	"cleo.com/testsupport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserService_ParseClinicalNote_ForWeightMetric(t *testing.T) {
	tests := []struct {
		desc           string
		clinicalNote   *model.ClinicalNote
		expectedMetric *model.HealthMetric
		expectedError  error
	}{
		{
			desc:           "clinical note with no weight metric provided",
			clinicalNote:   &model.ClinicalNote{Text: "lorem ipsum dolor sit amet"},
			expectedMetric: &model.HealthMetric{},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence weight of 75kg",
			clinicalNote:   &model.ClinicalNote{Text: "patient has provided a weight of 75 kg"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: wt of 75kg",
			clinicalNote:   &model.ClinicalNote{Text: "patient has provided a wt of 75 kg"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: weighs 75kg",
			clinicalNote:   &model.ClinicalNote{Text: "patient weighs 75 kg"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: weight is 75kg",
			clinicalNote:   &model.ClinicalNote{Text: "patient weight is 75 kg"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note with weight metric repeated, first metric reported",
			clinicalNote:   &model.ClinicalNote{Text: "patient has provided a weight of 75 kilograms and a weight of 120 pounds"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: weight is 75kgs",
			clinicalNote:   &model.ClinicalNote{Text: "patient weight is 75 kgs"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: weight is 75kilogram",
			clinicalNote:   &model.ClinicalNote{Text: "patient weight is 75kilogram"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: weight is 75kilograms",
			clinicalNote:   &model.ClinicalNote{Text: "patient weight is 75kilograms"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: weight is 165.34lb",
			clinicalNote:   &model.ClinicalNote{Text: "patient weight is 165.34lb"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: weight is 165.34lbs",
			clinicalNote:   &model.ClinicalNote{Text: "patient weight is 165.34lbs"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: weight is 165.34 pound",
			clinicalNote:   &model.ClinicalNote{Text: "patient weight is 165.34 pound"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note has weight metric using keyword sequence: weight is 165.34 pounds",
			clinicalNote:   &model.ClinicalNote{Text: "patient weight is 165.34 pounds"},
			expectedMetric: &model.HealthMetric{Weight: "75 kg"},
		},
		{
			desc:           "clinical note with invalid weight metrics given in kgs, below min weight",
			clinicalNote:   &model.ClinicalNote{Text: "patient has provided a weight of 0.5kg"},
			expectedMetric: nil,
			expectedError:  errors.New("invalid weight of 0.5 kg"),
		},
		{
			desc:           "clinical note with invalid weight metrics given in kgs, exceeding max weight",
			clinicalNote:   &model.ClinicalNote{Text: "patient has provided a weight of 700kg"},
			expectedMetric: nil,
			expectedError:  errors.New("invalid weight of 700 kg"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {

			testService := service.NewParserService(testsupport.Logger())
			healthMetric, err := testService.ParseClinicalNote(tt.clinicalNote)

			if tt.expectedError == nil {
				require.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedError, err)
			}
			assert.Equal(t, tt.expectedMetric, healthMetric)

		})
	}
}

func TestParserService_ParseClinicalNote_ForHeightMetric(t *testing.T) {
	tests := []struct {
		desc           string
		clinicalNote   *model.ClinicalNote
		expectedMetric *model.HealthMetric
		expectedError  error
	}{
		{
			desc:           "clinical note with no height metric provided",
			clinicalNote:   &model.ClinicalNote{Text: "lorem ipsum dolor sit amet"},
			expectedMetric: &model.HealthMetric{},
			expectedError:  nil,
		},
		{
			desc:           "clinical note with height metric provided and matched on keyword of",
			clinicalNote:   &model.ClinicalNote{Text: "height of 100cm"},
			expectedMetric: &model.HealthMetric{Height: "100 cm"},
			expectedError:  nil,
		},
		{
			desc:           "clinical note with height metric provided and matched on keyword is",
			clinicalNote:   &model.ClinicalNote{Text: "height is 100cm"},
			expectedMetric: &model.HealthMetric{Height: "100 cm"},
			expectedError:  nil,
		},
		{
			desc:           "clinical note with height metric provided and matched on keyword at",
			clinicalNote:   &model.ClinicalNote{Text: "height at 100cm"},
			expectedMetric: &model.HealthMetric{Height: "100 cm"},
			expectedError:  nil,
		},
		{
			desc:           "clinical note with height metric provided and not matched on any keyword of/at/is",
			clinicalNote:   &model.ClinicalNote{Text: "height approximately 100cm"},
			expectedMetric: &model.HealthMetric{Height: ""},
			expectedError:  nil,
		},
		{
			desc:           "clinical note with invalid height metrics given in feet",
			clinicalNote:   &model.ClinicalNote{Text: "patient has provided a height of 75feet"},
			expectedMetric: nil,
			expectedError:  errors.New("invalid height of 2286 cm"),
		},
		{
			desc:           "clinical note with invalid height metrics given in feet",
			clinicalNote:   &model.ClinicalNote{Text: "patient has provided a height of 75feet"},
			expectedMetric: nil,
			expectedError:  errors.New("invalid height of 2286 cm"),
		},
		{
			desc:         "clinical note with height metric repeated",
			clinicalNote: &model.ClinicalNote{Text: "patient has provided a height of 75 inches and a height of 120 pounds"},
			expectedMetric: &model.HealthMetric{
				Height: "190.5 cm",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {

			testService := service.NewParserService(testsupport.Logger())
			healthMetric, err := testService.ParseClinicalNote(tt.clinicalNote)

			if tt.expectedError == nil {
				require.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedError, err)
			}
			assert.Equal(t, tt.expectedMetric, healthMetric)

		})
	}

}
