package location

import "fmt"

// Airtable field names
const (
	FieldName = "Name"
	FieldSlug = "Slug"
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

// Location represents a physical place served by the API.
type Location struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}

// ToAirtableFields converts a Location to Airtable fields format.
func (l *Location) ToAirtableFields() map[string]interface{} {
	return map[string]interface{}{
		FieldName: l.Name,
		FieldSlug: l.Slug,
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

	return &Location{
		ID:     id,
		Name:   getStringField(fields, FieldName),
		Slug:   getStringField(fields, FieldSlug),
	}, nil
}
