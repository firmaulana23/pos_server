-- Migration: Add menu_item_id to add_ons table to make add-ons dependent on menu items
-- Created: 2025-07-07
-- Database: PostgreSQL

-- Add menu_item_id column to add_ons table
ALTER TABLE add_ons ADD COLUMN menu_item_id INTEGER;

-- Add foreign key constraint (optional, but recommended for data integrity)
ALTER TABLE add_ons ADD CONSTRAINT fk_add_ons_menu_item 
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE CASCADE;

-- Add index for performance
CREATE INDEX idx_add_ons_menu_item_id ON add_ons(menu_item_id);

-- Add comment to the column
COMMENT ON COLUMN add_ons.menu_item_id IS 'Foreign key to menu_items table. NULL means global add-on available for all items';

-- Update existing add-ons to be global (menu_item_id = NULL) - no change needed as NULL is default
-- Existing add-ons will remain available for all menu items until explicitly assigned to specific items
