                        Table "public.documents"
   Column    |           Type           | Collation | Nullable | Default 
-------------+--------------------------+-----------+----------+---------
 id          | text                     |           | not null | 
 name        | text                     |           | not null | 
 type        | text                     |           | not null | 
 version     | text                     |           | not null | 
 tags        | character varying[]      |           |          | 
 file_path   | text                     |           | not null | 
 file_size   | bigint                   |           | not null | 
 status      | text                     |           | not null | 
 description | text                     |           |          | 
 library     | text                     |           |          | 
 content     | text                     |           |          | 
 created_at  | timestamp with time zone |           |          | 
 updated_at  | timestamp with time zone |           |          | 
Indexes:
    "documents_pkey" PRIMARY KEY, btree (id)
    "idx_documents_library" btree (library)
    "idx_documents_name" btree (name)
    "idx_documents_status" btree (status)
    "idx_documents_type" btree (type)
    "idx_documents_version" btree (version)


                    Table "public.document_versions"
   Column    |           Type           | Collation | Nullable | Default 
-------------+--------------------------+-----------+----------+---------
 id          | text                     |           | not null | 
 document_id | text                     |           | not null | 
 version     | text                     |           | not null | 
 file_path   | text                     |           | not null | 
 file_size   | bigint                   |           | not null | 
 status      | text                     |           | not null | 
 description | text                     |           |          | 
 content     | text                     |           |          | 
 created_at  | timestamp with time zone |           |          | 
 updated_at  | timestamp with time zone |           |          | 
Indexes:
    "document_versions_pkey" PRIMARY KEY, btree (id)
    "idx_document_version_unique" UNIQUE CONSTRAINT, btree (document_id, version)
    "idx_document_versions_document_id" btree (document_id)
    "idx_document_versions_version" btree (version)

                    Table "public.document_metadata"
   Column    |           Type           | Collation | Nullable | Default 
-------------+--------------------------+-----------+----------+---------
 id          | text                     |           | not null | 
 document_id | text                     |           | not null | 
 metadata    | jsonb                    |           |          | 
 created_at  | timestamp with time zone |           |          | 
 updated_at  | timestamp with time zone |           |          | 
Indexes:
    "document_metadata_pkey" PRIMARY KEY, btree (id)
    "idx_document_metadata_document_id" btree (document_id)