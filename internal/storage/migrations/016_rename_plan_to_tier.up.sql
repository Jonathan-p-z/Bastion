DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'subscriptions' AND column_name = 'plan'
  ) THEN
    ALTER TABLE subscriptions RENAME COLUMN plan TO tier;
  END IF;
END $$;
