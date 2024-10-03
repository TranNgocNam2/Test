ALTER table sessions
    DROP CONSTRAINT IF EXISTS fk_session_subject;
ALTER table materials
    DROP CONSTRAINT IF EXISTS fk_material_session;
ALTER table transcripts
    DROP CONSTRAINT IF EXISTS fk_transcript_subject;

DROP table IF EXISTS sessions CASCADE;
DROP table IF EXISTS materials CASCADE;
DROP table IF EXISTS transcripts CASCADE;