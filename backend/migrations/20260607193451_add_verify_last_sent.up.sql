ALTER TABLE public.verification_codes
ADD COLUMN last_sent_at TIMESTAMP NOT NULL DEFAULT NOW();