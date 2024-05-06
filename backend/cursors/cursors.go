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
	URL     string   `json:"url"`
	PosX    int      `json:"posX"`
	PosY    int      `json:"posY"`
}

type GetCursors struct {
	Cursors []*Cursor `json:"cursors"`
}

type GetCursorsParams struct {
	URL string
}

// Cursors returns all cursors for a given path.
//
//encore:api public method=GET path=/cursors
func Cursors(ctx context.Context, p *GetCursorsParams) (GetCursors, error) {
	if p.URL == "" {
		return GetCursors{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "specify url in url parameters",
		}
	}

	cursors, err := getCursorsByURLFromDB(ctx, p.URL)
	if err != nil {
		rlog.Error("failed to retrieve cursors", "error", err)
		return GetCursors{}, &errs.Error{
			Code:    errs.Internal,
			Message: "failed to retrieve cursors",
		}
	}
	return GetCursors{Cursors: cursors}, nil
}
