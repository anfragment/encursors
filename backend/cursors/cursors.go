package cursors

import (
	"context"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type CursorOS int

const (
	CursorOSMacOS CursorOS = iota
	CursorOSWindows
	CursorOSLinux
)

type Cursor struct {
	Id      string   `json:"id"`
	Country string   `json:"country"`
	OS      CursorOS `json:"os"`
	Path    string   `json:"path"`
	PosX    int      `json:"posX"`
	PosY    int      `json:"posY"`
}

type GetCursors struct {
	Cursors []*Cursor `json:"cursors"`
}

type GetCursorsParams struct {
	Path string
}

// Cursors returns all cursors for a given path.
//
//encore:api public method=GET path=/cursors
func Cursors(ctx context.Context, p *GetCursorsParams) (GetCursors, error) {
	if p.Path == "" {
		return GetCursors{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "specify path in url parameters",
		}
	}

	cursors, err := getCursorsByPathFromDB(ctx, p.Path)
	if err != nil {
		rlog.Error("failed to retrieve cursors", "error", err)
		return GetCursors{}, &errs.Error{
			Code:    errs.Internal,
			Message: "failed to retrieve cursors",
		}
	}
	return GetCursors{Cursors: cursors}, nil
}
