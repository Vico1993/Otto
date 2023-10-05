-- Add last_time_parsed column
ALTER TABLE "chats" ADD COLUMN "last_time_parsed" timestamp without time zone;

-- Update all chats to not received notifications from old articles
UPDATE chats set last_time_parsed=NOW()