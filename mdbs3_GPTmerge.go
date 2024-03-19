package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// AWS-S3 Konfiguration
	accessKey := "owncloud"
	secretKey := "i6lfi2rnaj4rfi3eoudm3egolr5k1x68"
	endpoint := "http://128.140.86.10:8000"
	region := "us-east-1"

	// Setzen Sie die AWS-Region und benutzerdefinierte HTTP-Transportkonfiguration.
	cfg := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, "")).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithS3ForcePathStyle(true)

	// Erstelle eine neue AWS-Session
	sess, err := session.NewSession(cfg)
	if err != nil {
		panic(err)
	}

	// Erstelle einen S3-Client
	svc := s3.New(sess)

	// Verbindung zur MySQL-Datenbank herstellen
	db, err := sql.Open("mysql", "owncloud:ngq0ckid0r@tcp(128.140.86.10:13306)/owncloud")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verbindung zur MySQL-Datenbank prüfen
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

	// Ergebnisse in der Konsole und in eine Datei schreiben
	outputPath := "Ausgaben_mdb_s3/combined_data.txt"
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Fprintln(file, "S3 Bucket Objekte:")
	fmt.Fprintln(file, "-----------------------------")
	// S3 Bucket Objekte abrufen und in die Datei schreiben
	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String("oc-primary"), // Bucket-Namen angeben
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range resp.Contents {
		metaDataInput := &s3.HeadObjectInput{
			Bucket: aws.String("oc-primary"),
			Key:    item.Key,
		}
		metaDataOutput, err := svc.HeadObject(metaDataInput)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(file, "Name: %s\n", *item.Key)
		fmt.Fprintf(file, "Größe: %d\n", *item.Size)
		fmt.Fprintf(file, "AcceptRanges: %s\n", *metaDataOutput.AcceptRanges)
		fmt.Fprintf(file, "ContentLength: %d\n", *metaDataOutput.ContentLength)
		fmt.Fprintf(file, "Typ: %s\n", *metaDataOutput.ContentType)
		fmt.Fprintf(file, "ETag: %s\n", *metaDataOutput.ETag)
		fmt.Fprintf(file, "Zeitstempel: %v\n", *metaDataOutput.LastModified)
		fmt.Fprintln(file)
	}

	fmt.Fprintln(file, "\nMySQL Datenbank Objekte:")
	fmt.Fprintln(file, "-----------------------------")
	// MySQL Datenbank Objekte abrufen und in die Datei schreiben
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
			fmt.Println(output) // Ausgabe in der Konsole
			fmt.Fprintln(file, output)
		}
		fmt.Println() // Zeilenumbruch in der Konsole
		fmt.Fprintln(file) // Zeilenumbruch in der Datei
	}

	fmt.Printf("Ausgabe wurde in %s gespeichert.\n", outputPath)
}  

