-- Initialize the POS System database
-- This file is run when the PostgreSQL container starts for the first time

-- Create the main database if it doesn't exist
-- (already created via POSTGRES_DB environment variable)

-- Set timezone
SET timezone = 'Asia/Jakarta';
