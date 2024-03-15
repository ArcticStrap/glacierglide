-- Add a new column `namespace` to the `pages` table
ALTER TABLE pages
ADD COLUMN namespace INT NOT NULL DEFAULT 0;
