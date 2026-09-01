# Aion Manifold

## Overview
Aion Manifold is a research and governance operating system focused on structured discovery, objective setting, evidence review, and approval workflows. It is designed for safe, auditable decision-making without enabling payment, banking, publishing, or destructive actions.

## Mission
Provide a transparent, constrained workflow for:
- research state capture
- objective definition
- findings synthesis
- approval requests
- governance checkpoints

## Safety Boundaries
The system must never perform or facilitate:
- payments
- banking operations
- publishing workflows
- destructive execution
- unrestricted external actions with material risk

## Core Capabilities
### 1. Research State
- Record context, assumptions, and active investigation status.
- Maintain a clear state of what is known, unknown, and pending.

### 2. Objective Setting
- Define goals, scope, constraints, and success criteria.
- Ensure objectives are explicit and bounded.

### 3. Findings
- Summarize evidence, conclusions, and unresolved questions.
- Distinguish confirmed facts from hypotheses.

### 4. Approval Requests
- Route decision items to human review before execution.
- Require explicit approval before any external or operational action.

## Operating Principles
- Safety before action
- Transparent state tracking
- Human approval gates for sensitive actions
- Minimal scope and explicit boundaries
- Clear auditability of reasoning and decisions

## Recommended Architecture
- Frontend: React + TypeScript + Vite
- State model: structured workflow state
- Governance layer: approval checkpoints and role-based review
- Execution boundary: blocked for unsafe domains

## Deployment Notes
- Repository target: aionmanifold-svg/aion-manifold
- Intended branch: main
- Build validation required before production deployment
- Deployment access requires repository write permissions for the connected integration

## Acceptance Criteria
The product is considered ready when:
- the project builds successfully
- the repo is accessible to the deployment integration
- the application is deployed from main
- the live app is verified in a production environment
- the safety boundaries remain enforced
