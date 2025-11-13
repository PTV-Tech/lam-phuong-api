package location

import "fmt"

// Status represents the status of a location.
type Status string

const (
	// StatusActive indicates the location is active and available.
	StatusActive Status = "Active"
	// StatusDisabled indicates the location is disabled and unavailable.
	StatusDisabled Status = "Disabled"
)

// String returns the string representation of the status.
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status is a valid value.
func (s Status) IsValid() bool {
	return s == StatusActive || s == StatusDisabled
}

// Airtable field names
const (
	FieldName   = "Name"
	FieldSlug   = "Slug"
	FieldStatus = "Status"
)

// Helper functions
func getStringField(fields map[string]interface{}, key string) string {
	if val, ok := fields[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getStatusFromFields extracts and validates the status from Airtable fields.
// Returns StatusActive as default if status is empty or not provided.
func getStatusFromFields(fields map[string]interface{}) (Status, error) {
	statusStr := getStringField(fields, FieldStatus)

	// Default to Active if status is empty
	if statusStr == "" {
		return StatusActive, nil
	}

	// Convert string to Status type
	status := Status(statusStr)

	// Validate status
	if !status.IsValid() {
		return "", fmt.Errorf("invalid status: %s (must be 'Active' or 'Disabled')", statusStr)
	}

	return status, nil
}

// Location represents a physical place served by the API.
type Location struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Status Status `json:"status"`
}

// ToAirtableFields converts a Location to Airtable fields format.
func (l *Location) ToAirtableFields() map[string]interface{} {
	return map[string]interface{}{
		FieldName:   l.Name,
		FieldSlug:   l.Slug,
		FieldStatus: l.Status.String(),
	}
}

// FromAirtable maps an Airtable record to a Location.
// The record should have an "id" field and a "fields" map.
func FromAirtable(record map[string]interface{}) (*Location, error) {
	// Safely extract ID
	id := ""
	if idVal, ok := record["id"]; ok {
		if idStr, ok := idVal.(string); ok {
			id = idStr
		}
	}

	// Safely extract fields
	fields, ok := record["fields"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid record: missing or invalid 'fields'")
	}

	// Extract and parse status
	status, err := getStatusFromFields(fields)
	if err != nil {
		return nil, err
	}

	return &Location{
		ID:     id,
		Name:   getStringField(fields, FieldName),
		Slug:   getStringField(fields, FieldSlug),
		Status: status,
	}, nil
}
