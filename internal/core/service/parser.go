package service

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"cleo.com/internal/core/domain/model"
	"github.com/sirupsen/logrus"
)

type Metric struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

var (
	weightRegex = regexp.MustCompile(`(?i)\b(?:weight|wt|weighs)\s*(?:of|is|at|:)?\s*(\d{1,4}(?:\.\d{1,2})?)\s*(kg|kgs|kilogram|kilograms|lb|lbs|pound|pounds)\b`)
	heightRegex = regexp.MustCompile(`(?i)\b(?:height|ht)\s*(?:of|is|at|:)?\s*(\d{1,5}(?:\.\d{1,2})?)\s*(cm|mm|m|metre|metres|meter|meters|ft|feet|foot|in|inch|inches)\b`)
)

// apply sensible medical ranges for weight and height
const (
	minWeightKg = 1.0
	maxWeightKg = 635.0 // world's heaviest recorded approx
	minHeightCm = 20.0
	maxHeightCm = 272.0 // tallest recorded approx
)

type ParserService struct {
	logger *logrus.Logger
}

func NewParserService(logger *logrus.Logger) *ParserService {
	return &ParserService{
		logger: logger,
	}
}

func (s ParserService) ParseClinicalNote(note *model.ClinicalNote) (*model.HealthMetric, error) {
	if note == nil {
		return nil, nil
	}
	response := &model.HealthMetric{}

	weightMetric, err := extractWeightMetric(note.Text)
	if err != nil {
		s.logger.Infof("error encountered extracting weight metric %s", err.Error())
		return nil, err
	}

	if weightMetric != nil {
		response.Weight = fmt.Sprintf("%g %s", weightMetric.Value, weightMetric.Unit)
	}

	heightMetric, err := extractHeightMetric(note.Text)
	if err != nil {
		s.logger.Infof("error encountered extracting height metric %s", err.Error())
		return nil, err
	}

	if heightMetric != nil {
		response.Height = fmt.Sprintf("%g %s", heightMetric.Value, heightMetric.Unit)
	}

	return response, nil
}

func extractWeightMetric(text string) (*Metric, error) {
	note := strings.ToLower(text)
	// TODO seek requirements if more than one instance of weight metric found.  For now return first occurrence of metric
	if m := weightRegex.FindStringSubmatch(note); m != nil {
		valStr := m[1]
		unitStr := normalizeWeightUnit(m[2])
		v, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse given weight of %s ", valStr)
		}
		if unitStr != "kg" {
			v = weightToKg(v, unitStr)
		}
		if isValidWeight(v) {
			return &Metric{Value: round(v, 2), Unit: "kg"}, nil
		} else {
			return nil, fmt.Errorf("invalid weight of %g kg", round(v, 2))
		}
	}

	return nil, nil
}
func extractHeightMetric(text string) (*Metric, error) {
	note := strings.ToLower(text)
	// TODO seek requirements if more than one instance of height metric found. For now return first occurrence of metric
	if m := heightRegex.FindStringSubmatch(note); m != nil {
		valStr := m[1]
		unitStr := normalizeHeightUnit(m[2])
		v, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse given height of %s ", valStr)
		}
		if unitStr != "cm" {
			v = heightToCm(v, unitStr)
		}
		if isValidHeight(v) {
			return &Metric{Value: round(v, 1), Unit: "cm"}, nil
		} else {
			return nil, fmt.Errorf("invalid height of %g cm", round(v, 2))
		}
	}

	return nil, nil
}

func normalizeWeightUnit(u string) string {
	u = strings.ToLower(u)
	switch u {
	case "kg", "kgs", "kilogram", "kilograms":
		return "kg"
	case "lb", "lbs", "pound", "pounds":
		return "lb"
	default:
		return u
	}
}

func normalizeHeightUnit(u string) string {
	u = strings.ToLower(u)
	switch u {
	case "cm":
		return "cm"
	case "mm":
		return "mm"
	case "m", "meter", "meters", "metre", "metres":
		return "m"
	case "ft", "feet", "foot":
		return "ft"
	case "in", "inch", "inches":
		return "in"
	default:
		return u
	}
}

func weightToKg(val float64, unit string) float64 {
	switch unit {
	case "kg":
		return val
	case "lb":
		return val * 0.45359237
	default:
		return val
	}
}

func heightToCm(val float64, unit string) float64 {
	switch unit {
	case "cm":
		return val
	case "mm":
		return val / 10.0
	case "m":
		return val * 100.0
	case "ft":
		// usually given as 5.9 for 5ft9 or provided as two numbers; here we treat as feet decimal
		return val * 30.48
	case "in":
		return val * 2.54
	default:
		return val
	}
}

func isValidWeight(v float64) bool {
	return v >= minWeightKg && v <= maxWeightKg && !math.IsNaN(v)
}

func isValidHeight(v float64) bool {
	return v >= minHeightCm && v <= maxHeightCm && !math.IsNaN(v)
}

func round(x float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(x*p) / p
}
