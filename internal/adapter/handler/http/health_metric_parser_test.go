package http_test

import (
	"encoding/json"
	"errors"
	netHTTP "net/http"
	"testing"

	"cleo.com/internal/adapter/handler/http"
	"cleo.com/internal/core/domain/model"
	"cleo.com/internal/core/port/mocks"
	"cleo.com/testsupport"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_CreateUser(t *testing.T) {

	testHealthMetric := model.HealthMetric{
		Weight: "175Kg",
		Height: "75cm",
	}

	validResponseBytes, err := json.Marshal(testHealthMetric)
	require.NoError(t, err)

	tests := []struct {
		desc          string
		parserService *mocks.HealthMetricParserServiceMock
		clinicalNote  *model.ClinicalNote

		expectedHttpStatus                 int
		expectedHttpBody                   string
		expectedParseClinicalNoteCallCount int
	}{
		{
			desc:          "empty payload returns invalid request",
			parserService: &mocks.HealthMetricParserServiceMock{},
			clinicalNote:  nil,

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"error":"invalid request"}`,
		},
		{
			desc:          "empty note text returns invalid request",
			parserService: &mocks.HealthMetricParserServiceMock{},
			clinicalNote:  &model.ClinicalNote{Text: ""},

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"error":"invalid request"}`,
		},
		{
			desc:          "note exceeding text limit returns invalid request",
			parserService: &mocks.HealthMetricParserServiceMock{},
			clinicalNote: &model.ClinicalNote{
				Text: gofakeit.Paragraph(1, 1, 200, ""),
			},

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"error":"invalid request"}`,
		},
		{
			desc: "returns internal service error if ParseService errors",
			parserService: &mocks.HealthMetricParserServiceMock{
				ParseClinicalNoteFunc: func(note *model.ClinicalNote) (*model.HealthMetric, error) {
					return nil, errors.New("test internal service error")
				},
			},
			clinicalNote: &model.ClinicalNote{
				Text: gofakeit.Paragraph(1, 1, 1, ""),
			},

			expectedHttpStatus:                 netHTTP.StatusInternalServerError,
			expectedHttpBody:                   `{"error":"internal server error"}`,
			expectedParseClinicalNoteCallCount: 1,
		},
		{
			desc: "note within text limit returns successfully",
			parserService: &mocks.HealthMetricParserServiceMock{
				ParseClinicalNoteFunc: func(note *model.ClinicalNote) (*model.HealthMetric, error) {
					return &testHealthMetric, nil
				},
			},
			clinicalNote: &model.ClinicalNote{
				Text: gofakeit.Paragraph(1, 1, 1, ""),
			},

			expectedHttpStatus:                 netHTTP.StatusCreated,
			expectedHttpBody:                   string(validResponseBytes),
			expectedParseClinicalNoteCallCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		testHandler := http.NewHealthMetricParserHandler(testsupport.Logger(), tt.parserService)
		c, w := testsupport.NewTestContext(tt.clinicalNote)

		t.Run(tt.desc, func(t *testing.T) {
			testHandler.Parse(c)
			assert.Equal(t, tt.expectedHttpStatus, w.Code)
			assert.JSONEq(t, tt.expectedHttpBody, w.Body.String())

			require.Equal(t, tt.expectedParseClinicalNoteCallCount, len(tt.parserService.ParseClinicalNoteCalls()))

		})
	}

}
