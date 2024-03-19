package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
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

	// Jetzt bist du mit dem S3-Bucket verbunden und kannst Operationen durchführen
	fmt.Println("Verbunden mit S3-Bucket")

	// Beispiel zum Auflisten der Objekte im Bucket
	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String("oc-primary"), // Bucket-Namen angeben
	})
	if err != nil {
		panic(err)
	}
	var sum int = 0

	// Ausgabe der Objekte im Bucket
	for _, item := range resp.Contents {
		sum = sum + 1
		// Objekt-Metadaten abrufen
		metaDataInput := &s3.HeadObjectInput{
			Bucket: aws.String("oc-primary"),
			Key:    item.Key,
		}
		metaDataOutput, err := svc.HeadObject(metaDataInput)
		if err != nil {
			panic(err)
		}

		// Objektnamen und Metadaten ausgeben
		fmt.Println("Name:", *item.Key)                              // Name
		fmt.Println("Größe:", *item.Size)                            // Größe
		fmt.Println("AcceptRanges:", *metaDataOutput.AcceptRanges)   // AcceptRanges
		fmt.Println("ContentLength:", *metaDataOutput.ContentLength) // ContentLength
		fmt.Println("Typ:", *metaDataOutput.ContentType)             // Typ
		fmt.Println("ETag:", *metaDataOutput.ETag)                   // ETag
		fmt.Println("Zeitstempel:", *metaDataOutput.LastModified)    // Letzte Änderung

		// Leerzeile einfügen, um die Ausgabe zu trennen
		fmt.Println()
	}
	fmt.Println(sum)
	// Abfrage, ob die Ausgabe in eine Datei gespeichert werden soll
	fmt.Print("Soll die Ausgabe in eine Datei gespeichert werden (j/n)? ")
	var saveOption string
	if _, err := fmt.Scan(&saveOption); err != nil {
		panic(err)
	}

	// Wenn die Antwort "j" oder "J" ist, die Ausgabe in die Datei schreiben
	if saveOption == "j" || saveOption == "J" {
		// Pfad zum Ausgabeordner festlegen
		outputPath := "Ausgaben_mdb_s3/s3_Abfrage_Daten.txt"
		// Ausgabe in eine Datei schreiben
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Schreibe die Ausgabe in die Datei
		fmt.Fprintf(file, "Verbunden mit S3-Bucket\n\n")

		for _, item := range resp.Contents {
			// Objekt-Metadaten abrufen
			metaDataInput := &s3.HeadObjectInput{
				Bucket: aws.String("oc-primary"),
				Key:    item.Key,
			}
			metaDataOutput, err := svc.HeadObject(metaDataInput)
			if err != nil {
				panic(err)
			}

			// Objektnamen und Metadaten in die Datei schreiben
			fmt.Fprintf(file, "Name: %s\n", *item.Key)
			fmt.Fprintf(file, "Größe: %d\n", *item.Size)
			fmt.Fprintf(file, "AcceptRanges: %s\n", *metaDataOutput.AcceptRanges)
			fmt.Fprintf(file, "ContentLength: %d\n", *metaDataOutput.ContentLength)
			fmt.Fprintf(file, "Typ: %s\n", *metaDataOutput.ContentType)
			fmt.Fprintf(file, "ETag: %s\n", *metaDataOutput.ETag)
			fmt.Fprintf(file, "Zeitstempel: %v\n", *metaDataOutput.LastModified)
			// Leerzeile einfügen, um die Ausgabe zu trennen
			fmt.Fprintln(file)
		}

		fmt.Printf("Ausgabe wurde in %s gespeichert.\n", outputPath)
	}

}
