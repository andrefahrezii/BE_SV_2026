-- Fix audit log action constraint to include 'thrash'
ALTER TABLE IF EXISTS sv_portal.audit_logs
    DROP CONSTRAINT IF EXISTS audit_logs_action_check;

ALTER TABLE IF EXISTS sv_portal.audit_logs
    ADD CONSTRAINT audit_logs_action_check
    CHECK (action IN ('create','update','delete','login','logout','thrash'));
