// Copyright (c) 2021, 2023 Tim van der Molen <tim@kariliq.nl>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package signal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tbvdm/sigtop/sqlcipher"
)

type Context struct {
	dir                        string
	db                         *sqlcipher.DB
	dbVersion                  int
	recipientsByConversationID map[string]*Recipient
	recipientsByPhone          map[string]*Recipient
	recipientsByUUID           map[string]*Recipient
}

func Open(dir string) (*Context, error) {
	dbFile := filepath.Join(dir, DatabaseFile)

	// SQLite/SQLCipher doesn't provide a useful error message if the
	// database doesn't exist or can't be read
	f, err := os.Open(dbFile)
	if err != nil {
		return nil, err
	}
	f.Close()

	db, err := sqlcipher.OpenFlags(dbFile, sqlcipher.OpenReadOnly)
	if err != nil {
		return nil, err
	}

	key, err := dbKey(dir)
	if err != nil {
		db.Close()
		return nil, err
	}

	if err := db.Key(key); err != nil {
		db.Close()
		return nil, err
	}

	// Verify key
	if err := db.Exec("SELECT count(*) FROM sqlite_master"); err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot verify key: %w", err)
	}

	dbVersion, err := databaseVersion(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	if dbVersion < 19 {
		db.Close()
		return nil, fmt.Errorf("database version %d not supported (yet)", dbVersion)
	}

	ctx := Context{
		dir:       dir,
		db:        db,
		dbVersion: dbVersion,
	}

	return &ctx, nil
}

func (c *Context) Close() {
	c.db.Close()
}

func dbKey(dir string) ([]byte, error) {
	configFile := filepath.Join(dir, ConfigFile)
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config struct {
		Key string `json:"key"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("cannot parse %s: %w", configFile, err)
	}

	// Write the key as an SQLite blob literal
	key := "x'" + config.Key + "'"

	return []byte(key), nil
}
