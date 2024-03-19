package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Verbindung zur Datenbank herstellen
	db, err := sql.Open("mysql", "owncloud:ngq0ckid0r@tcp(128.140.86.10:13306)/owncloud")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verbindung zur Datenbank prüfen
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Abfrage ausführen
	rows, err := db.Query("SELECT name, fileid, checksum, size FROM oc_filecache")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Ergebnisse durchgehen
	for rows.Next() {
		var name, fileid string
		var checksum sql.NullString // Verwenden Sie sql.NullString, um NULL-Werte zu behandeln
		var size int64

		// Werte scannen und in Variablen speichern
		if err := rows.Scan(&name, &fileid, &checksum, &size); err != nil {
			log.Fatal(err)
		}

		// Ergebnisse ausgeben
		if checksum.Valid {
			fmt.Printf("Name: %s\nFileID: %s\nChecksum: %s\nSize: %d\n\n", name, fileid, checksum.String, size)
		} else {
			fmt.Printf("Name: %s\nFileID: %s\nChecksum: NULL\nSize: %d\n\n", name, fileid, size)
		}
	}

	// Fehlerüberprüfung
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

