package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	// Ergebnisse ausgeben
	for rows.Next() {
		// Werte scannen und in das values-Slice speichern
		if err := rows.Scan(valuePointers...); err != nil {
			log.Fatal(err)
		}

		// Ergebnisse in der Konsole ausgeben
		for i, value := range values {
			// Konvertiere Byte-Slices zu Strings, falls die Spalte eine Zeichenfolge enthält
			var output string
			if bytes, ok := value.([]byte); ok {
				output = fmt.Sprintf("%s: %s", columns[i], string(bytes))
			} else {
				output = fmt.Sprintf("%s: %v", columns[i], value)
			}
			fmt.Println(output) // Ausgabe in der Konsole
		}
		fmt.Println() // Zeilenumbruch hinzufügen
	}

	// Abfrage, ob die Ausgabe in eine Datei gespeichert werden soll
	fmt.Print("Soll die Ausgabe in eine Datei gespeichert werden (j/n)? ")
	var saveOption string
	if _, err := fmt.Scan(&saveOption); err != nil {
		log.Fatal(err)
	}

	// Wenn die Antwort "j" oder "J" ist, die Ausgabe in die Datei schreiben
	if saveOption == "j" || saveOption == "J" {
		// Pfad zum Ausgabeordner festlegen
		outputPath := "Ausgaben_mdb_s3/mdb_Abfrage_Daten.txt"
		// Ausgabe in eine Datei schreiben
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Zurück zum Anfang des Ergebnis-Sets
		rows, err := db.Query("SELECT * FROM oc_filecache")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Ergebnisse in die Datei schreiben
		for rows.Next() {
			if err := rows.Scan(valuePointers...); err != nil {
				log.Fatal(err)
			}
			for i, value := range values {
				var output string
				if bytes, ok := value.([]byte); ok {
					output = fmt.Sprintf("%s: %s", columns[i], string(bytes))
				} else {
					output = fmt.Sprintf("%s: %v", columns[i], value)
				}
				if _, err := file.WriteString(output + "\n"); err != nil {
					log.Fatal(err)
				}
			}
			if _, err := file.WriteString("\n"); err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("Ausgabe wurde in %s gespeichert.\n", outputPath)
	}
}

