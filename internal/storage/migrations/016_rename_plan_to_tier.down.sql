DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'subscriptions' AND column_name = 'tier'
  ) THEN
    ALTER TABLE subscriptions RENAME COLUMN tier TO plan;
  END IF;
END $$;
