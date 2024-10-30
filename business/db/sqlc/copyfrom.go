// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: copyfrom.go

package sqlc

import (
	"context"
)

// iteratorForInsertMaterial implements pgx.CopyFromSource.
type iteratorForInsertMaterial struct {
	rows                 []InsertMaterialParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertMaterial) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertMaterial) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ID,
		r.rows[0].SessionID,
		r.rows[0].Index,
		r.rows[0].Type,
		r.rows[0].IsShared,
		r.rows[0].Name,
		r.rows[0].Data,
	}, nil
}

func (r iteratorForInsertMaterial) Err() error {
	return nil
}

func (q *Queries) InsertMaterial(ctx context.Context, arg []InsertMaterialParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"materials"}, []string{"id", "session_id", "index", "type", "is_shared", "name", "data"}, &iteratorForInsertMaterial{rows: arg})
}

// iteratorForInsertSubjectSkill implements pgx.CopyFromSource.
type iteratorForInsertSubjectSkill struct {
	rows                 []InsertSubjectSkillParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertSubjectSkill) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertSubjectSkill) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ID,
		r.rows[0].SubjectID,
		r.rows[0].SkillID,
	}, nil
}

func (r iteratorForInsertSubjectSkill) Err() error {
	return nil
}

func (q *Queries) InsertSubjectSkill(ctx context.Context, arg []InsertSubjectSkillParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"subject_skills"}, []string{"id", "subject_id", "skill_id"}, &iteratorForInsertSubjectSkill{rows: arg})
}

// iteratorForInsertTranscripts implements pgx.CopyFromSource.
type iteratorForInsertTranscripts struct {
	rows                 []InsertTranscriptsParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertTranscripts) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertTranscripts) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ID,
		r.rows[0].SubjectID,
		r.rows[0].Name,
		r.rows[0].Index,
		r.rows[0].MinGrade,
		r.rows[0].Weight,
	}, nil
}

func (r iteratorForInsertTranscripts) Err() error {
	return nil
}

func (q *Queries) InsertTranscripts(ctx context.Context, arg []InsertTranscriptsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"transcripts"}, []string{"id", "subject_id", "name", "index", "min_grade", "weight"}, &iteratorForInsertTranscripts{rows: arg})
}
