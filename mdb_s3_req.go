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

func main () {
	// Konfiguration der AWS-S3
	accessKey := "owncloud"
	secretKey := "i6lfi2rnaj4rfi3eoudm3egolr5k1x68"
	endpoit := "http://128.140.86.10:8000"
	region := "us-east-1"
