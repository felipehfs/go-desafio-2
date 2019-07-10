CREATE TABLE files
(
  id uuid NOT NULL,
  created_at timestamp without time zone NOT NULL DEFAULT now(),
  updated_at timestamp without time zone,
  items jsonb,
  CONSTRAINT file_pkey PRIMARY KEY (id)
)