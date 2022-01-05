package plantuml

import "fmt"

const (
	RelationTypeExtension RelationType = iota
	RelationTypeComposition
	RelationTypeAggregation
	RelationTypeAlias
	RelationTypeArrow
)

type (
	RelationType int

	relation struct {
		from         RelationTarget
		to           RelationTarget
		relationType RelationType
	}

	RelationTarget struct {
		Namespace string
		Name      string
	}
)

func NewRelationTarget(name string) RelationTarget {
	return RelationTarget{
		Name: name,
	}
}

func NewRelationTargetWithNamespace(namespace string, name string) RelationTarget {
	return RelationTarget{
		Namespace: namespace,
		Name:      name,
	}
}

func (r *RelationTarget) String() string {
	if r.Namespace == "" {
		return r.Name
	}
	return fmt.Sprintf("%s.%s", r.Namespace, r.Name)
}

func (r *relation) buildRelationType() string {
	switch r.relationType {
	case RelationTypeExtension:
		return `<|--`
	case RelationTypeComposition:
		return `*--`
	case RelationTypeAggregation:
		return `o--`
	case RelationTypeAlias:
		return `#..`
	case RelationTypeArrow:
		return `<--`
	default:
		return `--`
	}
}

func (r *relation) Write(builder *LineStringBuilder, indent int) {
	typ := r.buildRelationType()
	builder.WriteLineWithDepth(indent, fmt.Sprintf(
		`"%s" %s "%s"`,
		r.to.String(),
		typ,
		r.from.String(),
	))
}

func Relation(from, to RelationTarget, relationType RelationType) Element {
	return &relation{
		from:         from,
		to:           to,
		relationType: relationType,
	}
}
