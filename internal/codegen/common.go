package codegen

import "sort"

type AttrDef struct {
	Name       string
	Type       string
	IsNullable bool
	IsID       bool
}

// OrderAttributes sorts attributes with identifier first, then non-nullable, then nullable.
func OrderAttributes(identifier string, attrs []struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsNullable bool   `json:"isNullable"`
}) []AttrDef {
	out := make([]AttrDef, 0, len(attrs))
	for _, a := range attrs {
		if a.Name == "" {
			continue
		}
		out = append(out, AttrDef{
			Name:       a.Name,
			Type:       a.Type,
			IsNullable: a.IsNullable,
			IsID:       identifier != "" && a.Name == identifier,
		})
	}

	sort.SliceStable(out, func(i, j int) bool {
		wi := weight(out[i])
		wj := weight(out[j])
		if wi != wj {
			return wi < wj
		}
		return false
	})
	return out
}

func weight(a AttrDef) int {
	if a.IsID {
		return 0
	}
	if !a.IsNullable {
		return 1
	}
	return 2
}
