package main

import (
    "fmt"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    
)


func main() {
    accessKey := "owncloud"
    secretKey := "1s68uyue47mdleketrbi5z72a9rc2kuf"
  //  hostBase := "65.21.154.79:8000"
    endpoint := "http://65.21.154.79:8000"
    region := "us-east-1"

// Setzen Sie die AWS-Region und benutzerdefinierte HTTP-Transportkonfiguration.
    cfg := aws.NewConfig().
        WithCredentials(credentials.NewStaticCredentials(accessKey,secretKey, "")).
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
    Bucket: aws.String("oc-primary"),
})
if err != nil {
    panic(err)
}

for _, item := range resp.Contents {
    fmt.Println("Name:", *item.Key)
    fmt.Println("Größe:", *item.Size)
}



}



