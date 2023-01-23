package models

import "strconv"

type SmartSearchColor struct {
	HexColor         string  `json:"hex_color"`
	DifferenceMethod string  `json:"difference_method"`
	Threshold        float64 `json:"threshold"`
}

type SmartSearchParams struct {
	SourceId         string           `json:"source_id"`
	StartDateTimeStr string           `json:"start_date_time_str"`
	EndDateTimeStr   string           `json:"end_date_time_str"`
	PredClassName    string           `json:"pred_class_name"`
	Color            SmartSearchColor `json:"color"`
}

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type Hex string

func (h Hex) toRGB() (RGB, error) {
	return Hex2RGB(h)
}

func Hex2RGB(hex Hex) (RGB, error) {
	var rgb RGB
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}
