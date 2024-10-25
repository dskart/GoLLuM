package draw

import (
	"fmt"
	"io"
	"text/template"

	"github.com/dskart/gollum/ringchain"
)

const dotTemplate = `strict {{.GraphType}} {
{{range $k, $v := .Attributes}}
	{{$k}}="{{$v}}";
{{end}}
{{range $s := .Statements}}
	"{{.Source}}" {{if .Target}}{{$.EdgeOperator}} "{{.Target}}" [ {{range $k, $v := .EdgeAttributes}}{{$k}}="{{$v}}", {{end}} weight={{.EdgeWeight}} ]{{else}}[ {{range $k, $v := .SourceAttributes}}{{$k}}="{{$v}}", {{end}} weight={{.SourceWeight}} ]{{end}};
{{end}}
}
`

type description struct {
	GraphType    string
	Attributes   map[string]string
	EdgeOperator string
	Statements   []statement
}

type statement struct {
	Source           interface{}
	Target           interface{}
	SourceWeight     int
	SourceAttributes map[string]string
	EdgeWeight       int
	EdgeAttributes   map[string]string
}

func DOT(g *ringchain.Graph, w io.Writer, options ...func(*description)) error {
	desc, err := generateDOT(g, options...)
	if err != nil {
		return fmt.Errorf("failed to generate DOT description: %w", err)
	}

	return renderDOT(w, desc)
}

func GraphAttribute(key, value string) func(*description) {
	return func(d *description) {
		d.Attributes[key] = value
	}
}

func generateDOT(g *ringchain.Graph, options ...func(*description)) (description, error) {
	desc := description{
		GraphType:    "graph",
		Attributes:   make(map[string]string),
		EdgeOperator: "--",
		Statements:   make([]statement, 0),
	}

	for _, option := range options {
		option(&desc)
	}

	desc.GraphType = "digraph"
	desc.EdgeOperator = "->"

	adjacencyMap, err := g.SuccessorMap()
	if err != nil {
		return desc, err
	}

	for node, adjacencies := range adjacencyMap {
		_, err := g.Node(node)
		if err != nil {
			return desc, err
		}

		stmt := statement{
			Source: node,
		}
		desc.Statements = append(desc.Statements, stmt)

		for adjacency := range adjacencies {
			stmt := statement{
				Source: node,
				Target: adjacency,
			}
			desc.Statements = append(desc.Statements, stmt)
		}
	}

	return desc, nil
}

func renderDOT(w io.Writer, d description) error {
	tpl, err := template.New("dotTemplate").Parse(dotTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tpl.Execute(w, d)
}
