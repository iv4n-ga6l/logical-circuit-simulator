package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
)

var inputSchemaLoader = gojsonschema.NewReferenceLoader("file://schemas/input_schema.json")

type ErrorResponse struct {
	Error           string      `json:"error"`
	Details         string      `json:"details"`
	InvalidComponent interface{} `json:"invalidComponent,omitempty"`
}

func SimulationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Validate input against schema
	result, err := gojsonschema.Validate(inputSchemaLoader, gojsonschema.NewBytesLoader(body))
	if err != nil {
		http.Error(w, "Failed to validate input schema", http.StatusInternalServerError)
		return
	}

	if !result.Valid() {
		var validationErrors []string
		for _, err := range result.Errors() {
			validationErrors = append(validationErrors, err.String())
		}
		response := ErrorResponse{
			Error:   "Invalid input",
			Details: "Input does not conform to schema",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Placeholder for further processing (e.g., circuit simulation)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Input validated successfully"}`))
}