package gameoflife

import "strings"

// RuleType is an enumeration for different rule types.
type RuleType int

const (
	// ConwayRuleType represents the Conway's Game of Life rule.
	ConwayRuleType RuleType = iota
	// NoTopLeftNeighborRuleType represents a rule that checks for the absence of a top left neighbor.
	NoTopLeftNeighborRuleType
)

var ruleNameToType = map[string]RuleType{
	"conway":      ConwayRuleType,
	"no-top-left": NoTopLeftNeighborRuleType,
}

func RuleFactory(ruleType RuleType) Rule {
	switch ruleType {
	case ConwayRuleType:
		return ConwayRule{}
	case NoTopLeftNeighborRuleType:
		return NoTopLeftNeighborRule{}
	default:
		// Return a default rule if no valid type is provided
		return ConwayRule{}
	}
}

// ParseRulesFromString parses comma-separated rule names into []Rule.
func ParseRulesFromString(rulesString string) []Rule {
	rules := make([]Rule, 0)

	if len(rulesString) == 0 {
		return rules
	}

	for _, ruleString := range strings.Split(rulesString, ",") {
		ruleString = strings.ToLower(strings.TrimSpace(ruleString))
		if ruleName, ok := ruleNameToType[ruleString]; ok {
			rules = append(rules, RuleFactory(ruleName))
		}
	}

	return rules
}

// AvailableRuleNames returns all valid rule names for CLI/help.
func AvailableRuleNames() []string {
	keys := make([]string, 0, len(ruleNameToType))
	for k := range ruleNameToType {
		keys = append(keys, k)
	}
	return keys
}
