//Copyright 2013 GoGraphviz Authors
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package gographviz

import (
	"github.com/awalterschulze/gographviz/ast"
)

// NewAnalysedGraph creates a Graph structure by analysing an Abstract Syntax Tree representing a parsed graph.
func NewAnalysedGraph(graph *ast.Graph) (*Graph, error) {
	g := NewGraph()
	if err := Analyse(graph, g); err != nil {
		return nil, err
	}
	return g, nil
}

// Analyse analyses an Abstract Syntax Tree representing a parsed graph into a newly created graph structure Interface.
func Analyse(graph *ast.Graph, g Interface) error {
	gerr := newErrCatcher(g)
	graph.Walk(&graphVisitor{gerr})
	return gerr.getError()
}

type nilVisitor struct {
}

func (w *nilVisitor) Visit(v ast.Elem) ast.Visitor {
	return w
}

type graphVisitor struct {
	g errInterface
}

func (w *graphVisitor) Visit(v ast.Elem) ast.Visitor {
	graph, ok := v.(*ast.Graph)
	if !ok {
		return w
	}
	w.g.SetStrict(graph.Strict)
	w.g.SetDir(graph.Type == ast.DIGRAPH)
	graphName := graph.ID.String()
	w.g.SetName(graphName)
	return newStmtVisitor(w.g, graphName)
}

func newStmtVisitor(g errInterface, graphName string) *stmtVisitor {
	return &stmtVisitor{g, graphName, make(map[string]string), make(map[string]string), make(map[string]string)}
}

type stmtVisitor struct {
	g                 errInterface
	graphName         string
	currentNodeAttrs  map[string]string
	currentEdgeAttrs  map[string]string
	currentGraphAttrs map[string]string
}

func (w *stmtVisitor) Visit(v ast.Elem) ast.Visitor {
	switch s := v.(type) {
	case ast.NodeStmt:
		return w.nodeStmt(s)
	case ast.EdgeStmt:
		return w.edgeStmt(s)
	case ast.NodeAttrs:
		return w.nodeAttrs(s)
	case ast.EdgeAttrs:
		return w.edgeAttrs(s)
	case ast.GraphAttrs:
		return w.graphAttrs(s)
	case *ast.SubGraph:
		return w.subGraph(s)
	case *ast.Attr:
		return w.attr(s)
	case ast.AttrList:
		return &nilVisitor{}
	default:
		//fmt.Fprintf(os.Stderr, "unknown stmt %T\n", v)
	}
	return w
}

func ammend(attrs map[string]string, add map[string]string) map[string]string {
	for key, value := range add {
		if _, ok := attrs[key]; !ok {
			attrs[key] = value
		}
	}
	return attrs
}

func overwrite(attrs map[string]string, overwrite map[string]string) map[string]string {
	for key, value := range overwrite {
		attrs[key] = value
	}
	return attrs
}

func (w *stmtVisitor) nodeStmt(stmt ast.NodeStmt) ast.Visitor {
	attrs := ammend(stmt.Attrs.GetMap(), w.currentNodeAttrs)
	w.g.AddNode(w.graphName, stmt.NodeID.String(), attrs)
	return &nilVisitor{}
}

func (w *stmtVisitor) edgeStmt(stmt ast.EdgeStmt) ast.Visitor {
	attrs := stmt.Attrs.GetMap()
	attrs = ammend(attrs, w.currentEdgeAttrs)
	src := stmt.Source.GetID()
	srcName := src.String()
	if stmt.Source.IsNode() {
		w.g.AddNode(w.graphName, srcName, w.currentNodeAttrs)
	}
	srcPort := stmt.Source.GetPort()
	for i := range stmt.EdgeRHS {
		directed := bool(stmt.EdgeRHS[i].Op)
		dst := stmt.EdgeRHS[i].Destination.GetID()
		dstName := dst.String()
		if stmt.EdgeRHS[i].Destination.IsNode() {
			w.g.AddNode(w.graphName, dstName, w.currentNodeAttrs)
		}
		dstPort := stmt.EdgeRHS[i].Destination.GetPort()
		w.g.AddPortEdge(srcName, srcPort.String(), dstName, dstPort.String(), directed, attrs)
		src = dst
		srcPort = dstPort
		srcName = dstName
	}
	return w
}

func (w *stmtVisitor) nodeAttrs(stmt ast.NodeAttrs) ast.Visitor {
	w.currentNodeAttrs = overwrite(w.currentNodeAttrs, ast.AttrList(stmt).GetMap())
	return &nilVisitor{}
}

func (w *stmtVisitor) edgeAttrs(stmt ast.EdgeAttrs) ast.Visitor {
	w.currentEdgeAttrs = overwrite(w.currentEdgeAttrs, ast.AttrList(stmt).GetMap())
	return &nilVisitor{}
}

func (w *stmtVisitor) graphAttrs(stmt ast.GraphAttrs) ast.Visitor {
	attrs := ast.AttrList(stmt).GetMap()
	for key, value := range attrs {
		w.g.AddAttr(w.graphName, key, value)
	}
	w.currentGraphAttrs = overwrite(w.currentGraphAttrs, attrs)
	return &nilVisitor{}
}

func (w *stmtVisitor) subGraph(stmt *ast.SubGraph) ast.Visitor {
	subGraphName := stmt.ID.String()
	w.g.AddSubGraph(w.graphName, subGraphName, w.currentGraphAttrs)
	return newStmtVisitor(w.g, subGraphName)
}

func (w *stmtVisitor) attr(stmt *ast.Attr) ast.Visitor {
	w.g.AddAttr(w.graphName, stmt.Field.String(), stmt.Value.String())
	return w
}
