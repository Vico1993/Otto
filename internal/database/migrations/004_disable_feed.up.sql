-- Add disabled column
ALTER TABLE feeds ADD COLUMN "disabled" BOOLEAN NOT NULL DEFAULT FALSE;
