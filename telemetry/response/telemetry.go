package response

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

//Telemetry is a struct that represents a given query result's data point
type Telemetry struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Value     float64   `json:"value,omitempty"`
	Metric    string    `json:"metric,omitempty"`
	Fields    map[string]interface{}
}

func (t *Telemetry) UnmarshalJSON(js []byte) error {
	dec := json.NewDecoder(bytes.NewReader(js))
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch tok.(type) {
		case string:
			val, err := dec.Token()
			if err != nil {
				return err
			}

			if tok == "timestamp" {
				t.Timestamp, _ = time.Parse(time.RFC3339, val.(string))
			} else if tok == "value" {
				t.Value = val.(float64)
			} else if tok == "metric" {
				t.Metric = val.(string)
			} else {
				if t.Fields == nil {
					t.Fields = map[string]interface{}{}
				}

				t.Fields[tok.(string)] = val
			}
		default:
			continue
		}
	}
	return nil
}
