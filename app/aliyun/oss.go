package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strconv"
	"strings"
	"time"
)

type AliyunOss struct {
	Endpoint string
	AccessKeyId string
	AccessKeySecret string
	BucketName string //yourBucketName
	BucketDomain string //Bucket domain
}


//Get the list of bucket names
func (c AliyunOss) GetListBuckets() ([]string , error) {
	client, err := oss.New(c.Endpoint , c.AccessKeyId , c.AccessKeySecret)
	if err != nil {
		return nil,err
	}

	lsRes, err := client.ListBuckets()
	if err != nil {
		return nil,err
	}

	result := []string{}
	for _, bucket := range lsRes.Buckets {
		result = append(result , bucket.Name)
	}

	return result,nil
}


//Uoload local files
//LocalFileName:local file name
//objectName:oss file name
func (c AliyunOss) UploadFile(localFileName string , objectName string) (string , error) {
	// Create OSSClient instance
	client, err := oss.New(c.Endpoint , c.AccessKeyId , c.AccessKeySecret)
	if err != nil {
		return "",err
	}
	// Get storage space
	bucket, err := client.Bucket(c.BucketName)
	if err != nil {
		return "",err
	}

	//storage by date
	date := time.Now()
	year := date.Year()
	month := date.Month()
	day  := date.Day()
	objectName = strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + strconv.Itoa(day) + "/" + objectName

	//upload files
	err = bucket.PutObjectFromFile(objectName , localFileName)
	if err != nil {
		return "",err
	}

	return objectName , nil
}


/Ddelete files on OSS
func (c AliyunOss) DeleteFile(objectName string) error {
	// create OSSClient instances
	client, err := oss.New(c.Endpoint , c.AccessKeyId , c.AccessKeySecret)
	if err != nil {
		return err
	}
	// get storage space
	bucket, err := client.Bucket(c.BucketName)
	if err != nil {
		return err
	}

	// Delete individual files. objectName indicates that the full path including the file suffix needs to be specified when deleting an OSS file, such as abc/efg/123.jpg.
	// To delete a folder, set objectName to the corresponding folder name. If the folder is not empty, you need to delete all objects under the folder before deleting the folder.
	err = bucket.DeleteObject(objectName)
	if err != nil {
		return err
	}
	return nil
}


//Get file url link
func (c AliyunOss) GetObjectFileUrl(objectFile string) string {
	if strings.Index(c.BucketDomain, "http://") == -1 && strings.Index(c.BucketDomain, "https://") == -1 {
		return "http://" + c.BucketDomain + "/" +  objectFile
	} else {
		return c.BucketDomain + "/" +  objectFile
	}
}
