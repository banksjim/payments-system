\set VERBOSITY terse

SELECT conname, contype
FROM pg_constraint
WHERE conrelid = 'entries'::regclass
ORDER BY conname;
