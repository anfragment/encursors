package cursors

import (
	"context"

	"encore.dev/storage/sqldb"
)

var db = sqldb.NewDatabase("cursors", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func writeCursorToDB(ctx context.Context, cursor *Cursor) error {
	_, err := db.Exec(ctx, "INSERT INTO cursors (id, country, os, url, pos_x, pos_y) VALUES ($1, $2, $3, $4, $5, $6)",
		cursor.Id, cursor.Country, cursor.OS, cursor.URL, cursor.PosX, cursor.PosY)
	return err
}

func getCursorsByURLFromDB(ctx context.Context, url string) ([]*Cursor, error) {
	rows, err := db.Query(ctx, "SELECT id, country, os, url, pos_x, pos_y FROM cursors WHERE url = $1", url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cursors []*Cursor
	for rows.Next() {
		var cursor Cursor
		if err := rows.Scan(&cursor.Id, &cursor.Country, &cursor.OS, &cursor.URL, &cursor.PosX, &cursor.PosY); err != nil {
			return nil, err
		}
		cursors = append(cursors, &cursor)
	}
	return cursors, nil
}

func updateCursorInDB(ctx context.Context, cursor *Cursor) error {
	_, err := db.Exec(ctx, "UPDATE cursors SET pos_x = $1, pos_y = $2 WHERE id = $3", cursor.PosX, cursor.PosY, cursor.Id)
	return err
}

func deleteCursorFromDB(ctx context.Context, id string) error {
	_, err := db.Exec(ctx, "DELETE FROM cursors WHERE id = $1", id)
	return err
}
