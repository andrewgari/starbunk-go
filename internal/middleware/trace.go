package middleware

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// AuditVerdict is the outcome of a single auditor node in a trace.
type AuditVerdict int

const (
	VerdictPassed  AuditVerdict = iota // auditor returned true
	VerdictFailed                      // auditor returned false
	VerdictIgnored                     // never evaluated due to short-circuit
)

// AuditNode is one node in the trace tree produced by TraceAudit.
type AuditNode struct {
	Name     string
	Verdict  AuditVerdict
	Children []AuditNode
}

// Namer is an optional interface that auditors may implement to provide a
// human-readable name for trace output. Named() is the primary way to attach
// a display name to any auditor.
type Namer interface {
	AuditorName() string
}

// Named wraps an auditor with a display name used in TraceAudit output.
// It has no effect on Audit behaviour.
func Named(name string, inner MessageAuditor) MessageAuditor {
	return namedAuditor{name: name, inner: inner}
}

type namedAuditor struct {
	name  string
	inner MessageAuditor
}

func (n namedAuditor) Audit(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return n.inner.Audit(s, m)
}

func (n namedAuditor) AuditorName() string { return n.name }

// TraceAudit evaluates a with the same short-circuit semantics as a.Audit,
// and returns a full AuditNode tree showing which checks passed, failed, or
// were short-circuited (VerdictIgnored).
func TraceAudit(a MessageAuditor, s *discordgo.Session, m *discordgo.MessageCreate) (bool, AuditNode) {
	switch v := a.(type) {
	case allOfAuditor:
		node := AuditNode{Name: "AllOf"}
		for i, child := range v.auditors {
			ok, childNode := TraceAudit(child, s, m)
			node.Children = append(node.Children, childNode)
			if !ok {
				// mark remaining children as ignored
				for _, remaining := range v.auditors[i+1:] {
					node.Children = append(node.Children, ignoredNode(remaining))
				}
				node.Verdict = VerdictFailed
				return false, node
			}
		}
		node.Verdict = VerdictPassed
		return true, node

	case anyOfAuditor:
		node := AuditNode{Name: "AnyOf"}
		for i, child := range v.auditors {
			ok, childNode := TraceAudit(child, s, m)
			node.Children = append(node.Children, childNode)
			if ok {
				for _, remaining := range v.auditors[i+1:] {
					node.Children = append(node.Children, ignoredNode(remaining))
				}
				node.Verdict = VerdictPassed
				return true, node
			}
		}
		node.Verdict = VerdictFailed
		return false, node

	case notAuditor:
		ok, childNode := TraceAudit(v.inner, s, m)
		verdict := VerdictPassed
		if ok {
			verdict = VerdictFailed
		}
		return !ok, AuditNode{Name: "Not", Verdict: verdict, Children: []AuditNode{childNode}}

	case namedAuditor:
		ok, childNode := TraceAudit(v.inner, s, m)
		verdict := VerdictPassed
		if !ok {
			verdict = VerdictFailed
		}
		return ok, AuditNode{Name: v.name, Verdict: verdict, Children: []AuditNode{childNode}}

	default:
		ok := a.Audit(s, m)
		verdict := VerdictPassed
		if !ok {
			verdict = VerdictFailed
		}
		return ok, AuditNode{Name: nameOf(a), Verdict: verdict}
	}
}

// FormatTrace returns a multi-line human-readable representation of an
// AuditNode tree, suitable for embedding in test failure messages.
func FormatTrace(root AuditNode) string {
	var b strings.Builder
	writeNode(&b, root, 0)
	return b.String()
}

func writeNode(b *strings.Builder, node AuditNode, depth int) {
	indent := strings.Repeat("  ", depth)
	verdict := verdictLabel(node.Verdict)
	// fixed-width name column: 22 chars (padded or truncated)
	nameCol := fmt.Sprintf("%-22s", node.Name)
	fmt.Fprintf(b, "%s%s %s\n", indent, nameCol, verdict)
	for _, child := range node.Children {
		writeNode(b, child, depth+1)
	}
}

func verdictLabel(v AuditVerdict) string {
	switch v {
	case VerdictPassed:
		return "PASSED"
	case VerdictFailed:
		return "FAILED"
	case VerdictIgnored:
		return "IGNORED"
	default:
		return "UNKNOWN"
	}
}

// ignoredNode builds a full VerdictIgnored subtree for an auditor without
// calling Audit on anything. Used to mark short-circuited children.
func ignoredNode(a MessageAuditor) AuditNode {
	switch v := a.(type) {
	case allOfAuditor:
		children := make([]AuditNode, len(v.auditors))
		for i, child := range v.auditors {
			children[i] = ignoredNode(child)
		}
		return AuditNode{Name: "AllOf", Verdict: VerdictIgnored, Children: children}
	case anyOfAuditor:
		children := make([]AuditNode, len(v.auditors))
		for i, child := range v.auditors {
			children[i] = ignoredNode(child)
		}
		return AuditNode{Name: "AnyOf", Verdict: VerdictIgnored, Children: children}
	case notAuditor:
		return AuditNode{Name: "Not", Verdict: VerdictIgnored, Children: []AuditNode{ignoredNode(v.inner)}}
	case namedAuditor:
		return AuditNode{Name: v.name, Verdict: VerdictIgnored, Children: []AuditNode{ignoredNode(v.inner)}}
	default:
		return AuditNode{Name: nameOf(a), Verdict: VerdictIgnored}
	}
}

// nameOf returns a display name for a leaf auditor. It checks for the Namer
// interface first, then falls back to the struct type name with the "Auditor"
// suffix stripped (e.g. "notSelfAuditor" → "notSelf").
func nameOf(a MessageAuditor) string {
	if n, ok := a.(Namer); ok {
		return n.AuditorName()
	}
	t := reflect.TypeOf(a)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	name := t.Name()
	if name == "" {
		return fmt.Sprintf("%T", a)
	}
	return strings.TrimSuffix(name, "Auditor")
}
