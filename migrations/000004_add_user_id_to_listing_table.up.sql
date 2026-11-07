ALTER TABLE listings 
ADD COLUMN user_id UUID NOT NULL REFERENCES users(id) 
ON DELETE CASCADE;

CREATE index idx_listings_user_id ON listings(user_id);