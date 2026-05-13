package authority

import (
	"strings"
)

func approvalReason(env AuthorityEnvelope, action ObservedAction, ruleRef string, resolution map[string]string) string {
	// approvalReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, req := range env.ApprovalRequirements {

		if reason := approvalRequirementReason(req, action, ruleRef, resolution); reason != "" {
			return reason
		}
	}
	return ""
}

func approvalRequirementReason(req ApprovalRequirement, action ObservedAction, ruleRef string, resolution map[string]string) string {
	// approvalRequirementReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if !approvalRequirementApplies(req, action, ruleRef) {
		return ""
	}
	if strings.TrimSpace(req.ApprovalEvidenceRef) == "" {

		return "approval_evidence_missing"
	}
	return evidenceRefsReason([]string{req.ApprovalEvidenceRef}, resolution)
}

func approvalRequirementApplies(req ApprovalRequirement, action ObservedAction, ruleRef string) bool {
	return (req.EventType == "" || req.EventType == action.EventType) &&
		(req.TargetRuleRef == "" || req.TargetRuleRef == ruleRef)
}
