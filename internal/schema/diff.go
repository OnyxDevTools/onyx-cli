package schema

import (
	"sort"
	"strings"
)

// ComputeSchemaDiff compares API schema with local schema.
func ComputeSchemaDiff(apiSchema SchemaRevision, localSchema SchemaUpsertRequest) SchemaDiff {
	apiMap := mapByName(apiSchema.Tables)
	localMap := mapByName(localSchema.Tables)

	var newTables, removedTables []string
	var changedTables []SchemaTableDiff

	for name := range localMap {
		if _, ok := apiMap[name]; !ok {
			newTables = append(newTables, name)
		}
	}
	for name := range apiMap {
		if _, ok := localMap[name]; !ok {
			removedTables = append(removedTables, name)
		}
	}

	for name, local := range localMap {
		apiEntity, ok := apiMap[name]
		if !ok {
			continue
		}
		if diff := diffEntity(apiEntity, local); diff != nil {
			changedTables = append(changedTables, *diff)
		}
	}

	sort.Strings(newTables)
	sort.Strings(removedTables)
	sort.Slice(changedTables, func(i, j int) bool { return changedTables[i].Name < changedTables[j].Name })

	return SchemaDiff{NewTables: newTables, RemovedTables: removedTables, ChangedTables: changedTables}
}

func mapByName(items []SchemaTable) map[string]SchemaTable {
	m := make(map[string]SchemaTable)
	for _, it := range items {
		if it.Name == "" {
			continue
		}
		m[it.Name] = it
	}
	return m
}

func diffEntity(api SchemaEntity, local SchemaEntity) *SchemaTableDiff {
	diff := SchemaTableDiff{Name: local.Name}
	var hasChange bool

	if strings.TrimSpace(api.Partition) != strings.TrimSpace(local.Partition) {
		diff.Partition = &PartitionChange{From: strings.TrimSpace(api.Partition), To: strings.TrimSpace(local.Partition)}
		hasChange = true
	}

	if !identifiersEqual(api.Identifier, local.Identifier) {
		diff.Identifier = &IdentifierChange{From: api.Identifier, To: local.Identifier}
		hasChange = true
	}

	if attrs := diffAttributes(api.Attributes, local.Attributes); attrs != nil {
		diff.Attributes = attrs
		hasChange = true
	}
	if idx := diffIndexes(api.Indexes, local.Indexes); idx != nil {
		diff.Indexes = idx
		hasChange = true
	}
	if res := diffResolvers(api.Resolvers, local.Resolvers); res != nil {
		diff.Resolvers = res
		hasChange = true
	}
	if trg := diffTriggers(api.Triggers, local.Triggers); trg != nil {
		diff.Triggers = trg
		hasChange = true
	}

	if !hasChange {
		return nil
	}
	return &diff
}

func identifiersEqual(a, b *SchemaIdentifier) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Name == b.Name && a.Type == b.Type && a.Generator == b.Generator
}

func diffAttributes(apiAttrs, localAttrs []SchemaAttribute) *AttributeChanges {
	apiMap := attrMap(apiAttrs)
	localMap := attrMap(localAttrs)
	var added []SchemaAttribute
	var removed []string
	var changed []AttributeChangeDetail

	for name, local := range localMap {
		if api, ok := apiMap[name]; ok {
			apiNull := api.IsNullable
			locNull := local.IsNullable
			if api.Type != local.Type || apiNull != locNull {
				changed = append(changed, AttributeChangeDetail{
					Name: name,
					From: AttributeInfo{Type: api.Type, IsNullable: apiNull},
					To:   AttributeInfo{Type: local.Type, IsNullable: locNull},
				})
			}
		} else {
			added = append(added, local)
		}
	}
	for name := range apiMap {
		if _, ok := localMap[name]; !ok {
			removed = append(removed, name)
		}
	}

	sort.Slice(added, func(i, j int) bool { return added[i].Name < added[j].Name })
	sort.Strings(removed)
	sort.Slice(changed, func(i, j int) bool { return changed[i].Name < changed[j].Name })

	if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
		return nil
	}
	return &AttributeChanges{Added: added, Removed: removed, Changed: changed}
}

func attrMap(attrs []SchemaAttribute) map[string]SchemaAttribute {
	m := make(map[string]SchemaAttribute)
	for _, a := range attrs {
		if a.Name == "" {
			continue
		}
		m[a.Name] = a
	}
	return m
}

func diffIndexes(apiIdx, localIdx []SchemaIndex) *IndexChanges {
	apiMap := idxMap(apiIdx)
	localMap := idxMap(localIdx)
	var added []SchemaIndex
	var removed []string
	var changed []IndexChangeDetail

	for name, local := range localMap {
		if api, ok := apiMap[name]; ok {
			apiType := api.Type
			if apiType == "" {
				apiType = "DEFAULT"
			}
			localType := local.Type
			if localType == "" {
				localType = "DEFAULT"
			}
			apiScore := api.MinimumScore
			locScore := local.MinimumScore
			if apiType != localType || !scoresEqual(apiScore, locScore) {
				changed = append(changed, IndexChangeDetail{Name: name, From: api, To: local})
			}
		} else {
			added = append(added, local)
		}
	}
	for name := range apiMap {
		if _, ok := localMap[name]; !ok {
			removed = append(removed, name)
		}
	}

	sort.Slice(added, func(i, j int) bool { return added[i].Name < added[j].Name })
	sort.Strings(removed)
	sort.Slice(changed, func(i, j int) bool { return changed[i].Name < changed[j].Name })

	if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
		return nil
	}
	return &IndexChanges{Added: added, Removed: removed, Changed: changed}
}

func idxMap(list []SchemaIndex) map[string]SchemaIndex {
	m := make(map[string]SchemaIndex)
	for _, it := range list {
		if it.Name == "" {
			continue
		}
		m[it.Name] = it
	}
	return m
}

func scoresEqual(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func diffResolvers(apiRes, localRes []SchemaResolver) *ResolverChanges {
	apiMap := resMap(apiRes)
	localMap := resMap(localRes)
	var added []SchemaResolver
	var removed []string
	var changed []ResolverChangeDetail

	for name, local := range localMap {
		if api, ok := apiMap[name]; ok {
			if api.Resolver != local.Resolver {
				changed = append(changed, ResolverChangeDetail{Name: name, From: api, To: local})
			}
		} else {
			added = append(added, local)
		}
	}
	for name := range apiMap {
		if _, ok := localMap[name]; !ok {
			removed = append(removed, name)
		}
	}

	sort.Slice(added, func(i, j int) bool { return added[i].Name < added[j].Name })
	sort.Strings(removed)
	sort.Slice(changed, func(i, j int) bool { return changed[i].Name < changed[j].Name })

	if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
		return nil
	}
	return &ResolverChanges{Added: added, Removed: removed, Changed: changed}
}

func resMap(list []SchemaResolver) map[string]SchemaResolver {
	m := make(map[string]SchemaResolver)
	for _, it := range list {
		if it.Name == "" {
			continue
		}
		m[it.Name] = it
	}
	return m
}

func diffTriggers(apiTrg, localTrg []SchemaTrigger) *TriggerChanges {
	apiMap := trgMap(apiTrg)
	localMap := trgMap(localTrg)
	var added []SchemaTrigger
	var removed []string
	var changed []TriggerChangeDetail

	for name, local := range localMap {
		if api, ok := apiMap[name]; ok {
			if api.Event != local.Event || api.Trigger != local.Trigger {
				changed = append(changed, TriggerChangeDetail{Name: name, From: api, To: local})
			}
		} else {
			added = append(added, local)
		}
	}
	for name := range apiMap {
		if _, ok := localMap[name]; !ok {
			removed = append(removed, name)
		}
	}

	sort.Slice(added, func(i, j int) bool { return added[i].Name < added[j].Name })
	sort.Strings(removed)
	sort.Slice(changed, func(i, j int) bool { return changed[i].Name < changed[j].Name })

	if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
		return nil
	}
	return &TriggerChanges{Added: added, Removed: removed, Changed: changed}
}

func trgMap(list []SchemaTrigger) map[string]SchemaTrigger {
	m := make(map[string]SchemaTrigger)
	for _, it := range list {
		if it.Name == "" {
			continue
		}
		m[it.Name] = it
	}
	return m
}
