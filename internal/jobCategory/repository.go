package jobcategory

import (
	"context"
	"fmt"
	"log"
	"strings"

	"lam-phuong-api/internal/airtable"
)

// Repository defines behavior for storing and retrieving job categories.
type Repository interface {
	List() []JobCategory
	Create(ctx context.Context, jobCategory JobCategory) (JobCategory, error)
	DeleteBySlug(slug string) bool
}

// AirtableRepository implements Repository interface using Airtable as the data store
type AirtableRepository struct {
	airtableClient *airtable.Client
	airtableTable  string
}

// NewAirtableRepository creates a repository that uses Airtable as the data store
func NewAirtableRepository(airtableClient *airtable.Client, airtableTable string) *AirtableRepository {
	return &AirtableRepository{
		airtableClient: airtableClient,
		airtableTable:  airtableTable,
	}
}

// List returns all job categories from Airtable
func (r *AirtableRepository) List() []JobCategory {
	records, err := r.airtableClient.ListRecords(context.Background(), r.airtableTable, nil)
	if err != nil {
		log.Printf("Failed to list job categories from Airtable: %v", err)
		return []JobCategory{} // Return empty slice on error
	}

	jobCategories := make([]JobCategory, 0, len(records))
	for _, record := range records {
		jc, err := mapAirtableRecord(record)
		if err != nil {
			log.Printf("Skipping Airtable record due to mapping error: %v", err)
			continue
		}
		jobCategories = append(jobCategories, jc)
	}

	return jobCategories
}

// Create adds a new job category to Airtable
func (r *AirtableRepository) Create(ctx context.Context, jobCategory JobCategory) (JobCategory, error) {
	// Save to Airtable
	airtableFields := jobCategory.ToAirtableFieldsForCreate()
	log.Printf("Attempting to save job category to Airtable table: %s", r.airtableTable)
	airtableRecord, err := r.airtableClient.CreateRecord(ctx, r.airtableTable, airtableFields)
	if err != nil {
		log.Printf("Failed to save job category to Airtable: %v", err)
		log.Printf("Error details - Table: %s, Fields: %+v", r.airtableTable, airtableFields)
		return JobCategory{}, fmt.Errorf("failed to create job category in Airtable: %w", err)
	}

	// Update the created job category with Airtable ID
	jobCategory.ID = airtableRecord.ID
	log.Printf("Job category saved to Airtable successfully with ID: %s", airtableRecord.ID)
	return jobCategory, nil
}

// DeleteBySlug removes a job category by its slug from Airtable
func (r *AirtableRepository) DeleteBySlug(slug string) bool {
	filterValue := escapeAirtableFormulaValue(slug)
	params := &airtable.ListParams{
		FilterByFormula: fmt.Sprintf("{%s} = '%s'", FieldSlug, filterValue),
	}

	records, err := r.airtableClient.ListRecords(context.Background(), r.airtableTable, params)
	if err != nil {
		log.Printf("Failed to query Airtable for slug %s: %v", slug, err)
		return false
	}

	if len(records) == 0 {
		return false
	}

	ids := make([]string, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.ID)
	}

	if err := r.airtableClient.BulkDeleteRecords(context.Background(), r.airtableTable, ids); err != nil {
		log.Printf("Failed to delete Airtable records for slug %s: %v", slug, err)
		return false
	}

	return true
}

func mapAirtableRecord(record airtable.Record) (JobCategory, error) {
	return JobCategory{
		ID:   record.ID,
		Name: getStringField(record.Fields, FieldName),
		Slug: getStringField(record.Fields, FieldSlug),
	}, nil
}

func escapeAirtableFormulaValue(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

