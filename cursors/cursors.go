package cursors

import (
	"context"
	"errors"

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
	Cursors []*Cursor
}

type GetCursorsParams struct {
	Path string
}

// Cursors returns all cursors for a given path.
//
//encore:api public method=GET path=/cursors
func Cursors(ctx context.Context, p *GetCursorsParams) (GetCursors, error) {
	if p.Path == "" {
		return GetCursors{}, errors.New("specify path in url parameters")
	}

	cursors, err := getCursorsByPathFromDB(ctx, p.Path)
	if err != nil {
		rlog.Error("failed to retrieve cursors", "error", err)
		return GetCursors{}, errors.New("failed to retrieve cursors")
	}
	return GetCursors{Cursors: cursors}, nil
}
