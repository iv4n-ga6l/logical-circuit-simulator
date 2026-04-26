package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "iv4n-ga6l/logical-circuit-simulator/handlers"
)

func TestSimulationHandler(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid predefined circuit",
			input:          `{"circuitType": "halfAdder", "inputs": {"A": true, "B": false}}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Input validated successfully"}`,
		},
		{
			name:           "Valid custom circuit",
			input:          `{"gates": [{"type": "AND", "name": "G1", "inputs": ["A", "B"]}], "outputs": ["G1"], "inputs": {"A": true, "B": false}}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Input validated successfully"}`,
		},
		{
			name:           "Invalid input schema",
			input:          `{"invalidField": "value"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid input","details":"Input does not conform to schema"}`,
		},
		{
			name:           "Undefined gate reference",
			input:          `{"gates": [{"type": "AND", "name": "G1", "inputs": ["A", "X"]}], "outputs": ["G1"], "inputs": {"A": true, "B": false}}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Undefined input reference","details":"Gate 'G1' references undefined input 'X'","invalidComponent":"G1"}`,
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/simulate", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handlers.SimulationHandler(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			if res.StatusCode != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, res.StatusCode)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if string(body) != test.expectedBody {
				t.Errorf("Expected body %s, got %s", test.expectedBody, string(body))
			}
		})
	}
}