package governance

import (
	"context"
	"fmt"
)

type AuthorityPhase string

const (
	Manual     AuthorityPhase = "manual"
	Connector  AuthorityPhase = "connector"
	Autonomous AuthorityPhase = "autonomous"
)

type ValidationStatus string

const (
	Verified    ValidationStatus = "verified"
	Unverified  ValidationStatus = "unverified"
	Quarantined ValidationStatus = "quarantined"
)

type Action struct {
	Name       string
	Parameters map[string]string
}

type ConnectorTransition struct {
	AgentID          string
	Phase            AuthorityPhase
	Granted          map[string]bool
	HasHumanApproval bool
}

type GovernorEvent struct {
	AgentID          string
	ActionAttempted  string
	Phase            AuthorityPhase
	IsGranted        bool
	HasHumanToken    bool
	IAEComputedEmpty bool
	ResolvedStatus   ValidationStatus
	EvidenceEBPFHash string
}

type GovernorSink interface {
	PublishIAE(ctx context.Context, event GovernorEvent) error
}

// EvaluateInvariant is the deterministic enforcement boundary.
// Model/RAG assertions never enlarge Granted.
func (ct *ConnectorTransition) EvaluateInvariant(
	ctx context.Context,
	actualAction Action,
	evidenceEBPFHash string,
	governor GovernorSink,
) (ValidationStatus, error) {
	granted := ct.Granted[actualAction.Name]

	// Explicit human approval may authorize an otherwise-ungranted action,
	// but it is recorded separately from the static GRANTED set.
	if granted {
		event := GovernorEvent{
			AgentID:          ct.AgentID,
			ActionAttempted:  actualAction.Name,
			Phase:            ct.Phase,
			IsGranted:        true,
			HasHumanToken:    false,
			IAEComputedEmpty: true,
			ResolvedStatus:   Verified,
			EvidenceEBPFHash: evidenceEBPFHash,
		}
		if governor != nil {
			if err := governor.PublishIAE(ctx, event); err != nil {
				return Unverified, fmt.Errorf("governor publication failed: %w", err)
			}
		}
		return Verified, nil
	}

	if ct.HasHumanApproval {
		event := GovernorEvent{
			AgentID:          ct.AgentID,
			ActionAttempted:  actualAction.Name,
			Phase:            ct.Phase,
			IsGranted:        false,
			HasHumanToken:    true,
			IAEComputedEmpty: false,
			ResolvedStatus:   Verified,
			EvidenceEBPFHash: evidenceEBPFHash,
		}
		if governor != nil {
			if err := governor.PublishIAE(ctx, event); err != nil {
				return Unverified, fmt.Errorf("governor publication failed: %w", err)
			}
		}
		return Verified, nil
	}

	// IAE = ACTUAL \ GRANTED is non-empty.
	// No execution token may be issued for this path.
	event := GovernorEvent{
		AgentID:          ct.AgentID,
		ActionAttempted:  actualAction.Name,
		Phase:            ct.Phase,
		IsGranted:        false,
		HasHumanToken:    false,
		IAEComputedEmpty: false,
		ResolvedStatus:   Quarantined,
		EvidenceEBPFHash: evidenceEBPFHash,
	}
	if governor != nil {
		if err := governor.PublishIAE(ctx, event); err != nil {
			return Unverified, fmt.Errorf("governor publication failed: %w", err)
		}
	}

	return Quarantined, fmt.Errorf(
		"invariant violation (IAE detected): action %s is not in GRANTED set",
		actualAction.Name,
	)
}
