-- ASTRA authority-phase tracking
-- Migration: 20260922000000_add_connector_authority_phases.sql

DO $$ BEGIN
  CREATE TYPE agent_authority_phase AS ENUM ('manual', 'connector', 'autonomous');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
  CREATE TYPE validation_status AS ENUM ('verified', 'unverified', 'quarantined');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

ALTER TABLE agent_environments
  ADD COLUMN IF NOT EXISTS authority_phase agent_authority_phase NOT NULL DEFAULT 'manual',
  ADD COLUMN IF NOT EXISTS granted_capabilities TEXT[] NOT NULL DEFAULT '{}',
  ADD COLUMN IF NOT EXISTS last_iae_status validation_status NOT NULL DEFAULT 'verified',
  ADD COLUMN IF NOT EXISTS last_checked_at TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS authority_governance_ledger (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  agent_id TEXT NOT NULL,
  action_attempted TEXT NOT NULL,
  phase agent_authority_phase NOT NULL,
  is_granted BOOLEAN NOT NULL,
  has_human_token BOOLEAN NOT NULL DEFAULT FALSE,
  iae_computed_empty BOOLEAN NOT NULL,
  evidence_ebpf_hash TEXT,
  resolved_status validation_status NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_gov_ledger_agent
  ON authority_governance_ledger(agent_id);
