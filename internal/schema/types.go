package schema

// SchemaRevision represents a published or latest schema revision.
type SchemaRevision struct {
	DatabaseID string        `json:"databaseId"`
	Tables     []SchemaTable `json:"tables"`
	Entities   []SchemaTable `json:"entities,omitempty"` // backwards compatibility
	Meta       *SchemaMeta   `json:"meta,omitempty"`
}

type SchemaMeta struct {
	RevisionID  string `json:"revisionId,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	PublishedAt string `json:"publishedAt,omitempty"`
}

type SchemaTable struct {
	Name       string            `json:"name"`
	Partition  string            `json:"partition,omitempty"`
	Identifier *SchemaIdentifier `json:"identifier,omitempty"`
	Attributes []SchemaAttribute `json:"attributes,omitempty"`
	Indexes    []SchemaIndex     `json:"indexes,omitempty"`
	Resolvers  []SchemaResolver  `json:"resolvers,omitempty"`
	Triggers   []SchemaTrigger   `json:"triggers,omitempty"`
}

// Backwards compatibility alias (legacy name was SchemaEntity).
type SchemaEntity = SchemaTable

type SchemaIdentifier struct {
	Name      string `json:"name,omitempty"`
	Generator string `json:"generator,omitempty"`
	Type      string `json:"type,omitempty"`
}

type SchemaAttribute struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsNullable bool   `json:"isNullable,omitempty"`
}

type SchemaIndex struct {
	Name         string   `json:"name"`
	Type         string   `json:"type,omitempty"`
	MinimumScore *float64 `json:"minimumScore,omitempty"`
}

type SchemaResolver struct {
	Name     string `json:"name"`
	Resolver string `json:"resolver"`
}

type SchemaTrigger struct {
	Name    string `json:"name"`
	Event   string `json:"event"`
	Trigger string `json:"trigger"`
}

type SchemaUpsertRequest struct {
	DatabaseID string        `json:"databaseId,omitempty"`
	Tables     []SchemaTable `json:"tables"`
	Entities   []SchemaTable `json:"entities,omitempty"` // backwards compatibility
	Meta       *SchemaMeta   `json:"meta,omitempty"`
}

type SchemaValidationResult struct {
	Valid  *bool             `json:"valid,omitempty"`
	Errors []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
	Message string `json:"message"`
}

type SchemaDiff struct {
	NewTables     []string          `json:"newTables"`
	RemovedTables []string          `json:"removedTables"`
	ChangedTables []SchemaTableDiff `json:"changedTables"`
}

type SchemaTableDiff struct {
	Name       string            `json:"name"`
	Partition  *PartitionChange  `json:"partition,omitempty"`
	Identifier *IdentifierChange `json:"identifier,omitempty"`
	Attributes *AttributeChanges `json:"attributes,omitempty"`
	Indexes    *IndexChanges     `json:"indexes,omitempty"`
	Resolvers  *ResolverChanges  `json:"resolvers,omitempty"`
	Triggers   *TriggerChanges   `json:"triggers,omitempty"`
}

type PartitionChange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type IdentifierChange struct {
	From *SchemaIdentifier `json:"from,omitempty"`
	To   *SchemaIdentifier `json:"to,omitempty"`
}

type AttributeChanges struct {
	Added   []SchemaAttribute       `json:"added,omitempty"`
	Removed []string                `json:"removed,omitempty"`
	Changed []AttributeChangeDetail `json:"changed,omitempty"`
}

type AttributeChangeDetail struct {
	Name string        `json:"name"`
	From AttributeInfo `json:"from"`
	To   AttributeInfo `json:"to"`
}

type AttributeInfo struct {
	Type       string `json:"type"`
	IsNullable bool   `json:"isNullable"`
}

type IndexChanges struct {
	Added   []SchemaIndex       `json:"added,omitempty"`
	Removed []string            `json:"removed,omitempty"`
	Changed []IndexChangeDetail `json:"changed,omitempty"`
}

type IndexChangeDetail struct {
	Name string      `json:"name"`
	From SchemaIndex `json:"from"`
	To   SchemaIndex `json:"to"`
}

type ResolverChanges struct {
	Added   []SchemaResolver       `json:"added,omitempty"`
	Removed []string               `json:"removed,omitempty"`
	Changed []ResolverChangeDetail `json:"changed,omitempty"`
}

type ResolverChangeDetail struct {
	Name string         `json:"name"`
	From SchemaResolver `json:"from"`
	To   SchemaResolver `json:"to"`
}

type TriggerChanges struct {
	Added   []SchemaTrigger       `json:"added,omitempty"`
	Removed []string              `json:"removed,omitempty"`
	Changed []TriggerChangeDetail `json:"changed,omitempty"`
}

type TriggerChangeDetail struct {
	Name string        `json:"name"`
	From SchemaTrigger `json:"from"`
	To   SchemaTrigger `json:"to"`
}
