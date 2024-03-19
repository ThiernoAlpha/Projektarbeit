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
	rows, err := db.Query("SELECT * FROM oc_filecache")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Spaltennamen abrufen
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// Ein Slice von leeren Schnittstellen erstellen, um die Spaltenwerte zu speichern
	values := make([]interface{}, len(columns))
	// Ein Slice von Pointern auf die Elemente des values-Slices erstellen
	valuePointers := make([]interface{}, len(columns))
	for i := range values {
		valuePointers[i] = &values[i]
	}

	// Ergebnisse durchgehen
	for rows.Next() {
		// Werte scannen und in das values-Slice speichern
		if err := rows.Scan(valuePointers...); err != nil {
			log.Fatal(err)
		}

		// Ergebnisse ausgeben
		for i, value := range values {
			// Konvertiere Byte-Slices zu Strings, falls die Spalte eine Zeichenfolge enthält
			if bytes, ok := value.([]byte); ok {
				fmt.Printf("%s: %s\n", columns[i], string(bytes))
			} else {
				fmt.Printf("%s: %v\n", columns[i], value)
			}
		}
		fmt.Println() // Zeilenumbruch hinzufügen
	}

	// Fehlerüberprüfung
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
