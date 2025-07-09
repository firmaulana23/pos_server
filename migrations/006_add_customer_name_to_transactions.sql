-- Migration: Add customer_name column to transactions table
-- Date: 2025-07-08
-- Description: Add customer name field to track customer information for orders

ALTER TABLE transactions ADD COLUMN customer_name VARCHAR(255) DEFAULT '';

-- Add comment for documentation
COMMENT ON COLUMN transactions.customer_name IS 'Customer name for the order (optional)';
