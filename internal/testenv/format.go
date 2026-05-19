package testenv

import (
	"fmt"
	"strings"

	"github.com/andrewgari/starbunk-go/internal/middleware"
)

// FormatScenarioFailure produces a multi-line string describing a scenario
// result. Pass it as the optional annotation argument to Gomega's Expect so
// the full audit tree appears in test failure output automatically.
//
// Example usage:
//
//	result := harness.Run(msg)
//	Expect(result.AuditPassed).To(BeTrue(), testenv.FormatScenarioFailure("my scenario", result))
func FormatScenarioFailure(description string, result ScenarioResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Scenario: %s\n", description)
	fmt.Fprintf(&b, "Audit:    %s\n", verdictWord(result.AuditTrace.Verdict))
	fmt.Fprintf(&b, "\nAudit tree:\n")
	writeNode(&b, result.AuditTrace, 0)

	if len(result.Messages) == 0 {
		fmt.Fprintf(&b, "\nMessages sent: none\n")
	} else {
		fmt.Fprintf(&b, "\nMessages sent:\n")
		for i, msg := range result.Messages {
			if msg.Username != "" {
				fmt.Fprintf(&b, "  [%d] channel=%s identity=%s content=%q\n",
					i+1, msg.ChannelID, msg.Username, msg.Content)
			} else {
				fmt.Fprintf(&b, "  [%d] channel=%s content=%q\n",
					i+1, msg.ChannelID, msg.Content)
			}
		}
	}
	return b.String()
}

func verdictWord(v middleware.AuditVerdict) string {
	switch v {
	case middleware.VerdictPassed:
		return "PASSED"
	case middleware.VerdictFailed:
		return "FAILED"
	case middleware.VerdictIgnored:
		return "IGNORED"
	default:
		return "UNKNOWN"
	}
}

func writeNode(b *strings.Builder, node middleware.AuditNode, depth int) {
	indent := strings.Repeat("  ", depth)
	nameCol := fmt.Sprintf("%-22s", node.Name)
	fmt.Fprintf(b, "%s%s %s\n", indent, nameCol, verdictWord(node.Verdict))
	for _, child := range node.Children {
		writeNode(b, child, depth+1)
	}
}
