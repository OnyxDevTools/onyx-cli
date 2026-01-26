package schema

import (
	"bytes"
	"sort"

	"gopkg.in/yaml.v3"
)

// FormatSchemaDiff returns a YAML-style string similar to the TypeScript CLI.
func FormatSchemaDiff(diff SchemaDiff, filePath string) string {
	hasChanges := len(diff.NewTables) > 0 || len(diff.RemovedTables) > 0 || len(diff.ChangedTables) > 0
	if !hasChanges {
		if filePath == "" {
			filePath = "local schema"
		}
		return "No differences found between API schema and " + filePath + ".\n"
	}

	buf := &bytes.Buffer{}
	header := "# Schema diff"
	if filePath != "" {
		header = "# Diff between API schema and " + filePath
	}
	buf.WriteString(header)
	buf.WriteString("\n\n")

	// Build a minimal object, omitting empty sections for readability.
	pruned := map[string]any{}

	if len(diff.NewTables) > 0 {
		pruned["newTables"] = append([]string{}, diff.NewTables...)
	}
	if len(diff.RemovedTables) > 0 {
		pruned["removedTables"] = append([]string{}, diff.RemovedTables...)
	}

	// prune changed tables and include only when non-empty
	var changed []SchemaTableDiff
	for _, t := range diff.ChangedTables {
		prunedTable := SchemaTableDiff{Name: t.Name}
		if t.Partition != nil && t.Partition.From != t.Partition.To {
			prunedTable.Partition = t.Partition
		}
		if t.Identifier != nil && !identifiersEqual(t.Identifier.From, t.Identifier.To) {
			prunedTable.Identifier = t.Identifier
		}
		if t.Attributes != nil && (len(t.Attributes.Added) > 0 || len(t.Attributes.Removed) > 0 || len(t.Attributes.Changed) > 0) {
			prunedTable.Attributes = t.Attributes
		}
		if t.Indexes != nil && (len(t.Indexes.Added) > 0 || len(t.Indexes.Removed) > 0 || len(t.Indexes.Changed) > 0) {
			prunedTable.Indexes = t.Indexes
		}
		if t.Resolvers != nil && (len(t.Resolvers.Added) > 0 || len(t.Resolvers.Removed) > 0 || len(t.Resolvers.Changed) > 0) {
			prunedTable.Resolvers = t.Resolvers
		}
		if t.Triggers != nil && (len(t.Triggers.Added) > 0 || len(t.Triggers.Removed) > 0 || len(t.Triggers.Changed) > 0) {
			prunedTable.Triggers = t.Triggers
		}
		if prunedTable.Partition != nil || prunedTable.Identifier != nil || prunedTable.Attributes != nil || prunedTable.Indexes != nil || prunedTable.Resolvers != nil || prunedTable.Triggers != nil {
			changed = append(changed, prunedTable)
		}
	}
	if len(changed) > 0 {
		sort.Slice(changed, func(i, j int) bool { return changed[i].Name < changed[j].Name })
		pruned["changedTables"] = changed
	}

	enc := yaml.NewEncoder(buf)
	enc.SetIndent(2)
	if len(pruned) == 0 {
		_ = enc.Encode(map[string]any{"note": "no differences"})
	} else {
		_ = enc.Encode(pruned)
	}
	enc.Close()
	return buf.String()
}
