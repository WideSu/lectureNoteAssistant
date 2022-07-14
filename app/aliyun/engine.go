package aliyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

const (
	ALIYUN_CLOUND_REGION_CHA = 1 //China
	ALIYUN_CLOUND_REGION_INT = 2 //Overseas
)

//Alibaba Cloud Speech Recognition Engine
type AliyunClound struct {
	AccessKeyId     string
	AccessKeySecret string
	AppKey          string
	Region          int
}

//Alibaba Cloud recording file recognition result set
type AliyunAudioRecognitionResult struct {
	Text            string //text results
	TranslateText   string //translate text results
	ChannelId       int64  //Track ID
	BeginTime       int64  //The start time offset of the sentence, in milliseconds
	EndTime         int64  //The end time offset of the sentence, in milliseconds
	SilenceDuration int64  //The duration of silence between this sentence and the previous sentence, in seconds
	SpeechRate      int64  //The average speaking rate of this sentence, in words per minute
	EmotionValue    int64  //The emotional energy value is 1-10, the higher the value, the stronger the emotion
}

//Alibaba Cloud Recognition Word Dataset
type AliyunAudioWord struct {
	Word      string
	ChannelId int64
	BeginTime int64
	EndTime   int64
}

// Befine constants
const REGION_ID_REGION_INT string = "ap-southeast-1"
const ENDPOINT_NAME_REGION_INT string = "ap-southeast-1"
const PRODUCT_REGION_INT string = "nls-filetrans"
const DOMAIN_REGION_INT string = "filetrans.ap-southeast-1.aliyuncs.com"
const API_VERSION_REGION_INT string = "2019-08-23"

const REGION_ID_REGION_CHA string = "cn-shanghai"
const ENDPOINT_NAME_REGION_CHA string = "cn-shanghai"
const PRODUCT_REGION_CHA string = "nls-filetrans"
const DOMAIN_REGION_CHA string = "filetrans.cn-shanghai.aliyuncs.com"
const API_VERSION_REGION_CHA string = "2018-08-17"

const POST_REQUEST_ACTION string = "SubmitTask"
const GET_REQUEST_ACTION string = "GetTaskResult"

// request parameter key
const KEY_APP_KEY string = "appkey"
const KEY_FILE_LINK string = "file_link"
const KEY_VERSION string = "version"
const KEY_ENABLE_WORDS string = "enable_words"

//Whether to enable unified post-processing, the default value is false
const ENABLE_UNIFY_POST string = "enable_unify_post"

//Whether to open ITN, Chinese numbers will be converted to Arabic numbers for output, the default value is false
//When opening, you need to set version to "4.0", enable_unify_post must be true
const ENABLE_INVERSE_TEXT_NORMALIZATION string = "enable_inverse_text_normalization"

//If you want to enable the post-processing model, the default value is chinese, and you need to set the version to "4.0" when you open it.
//enable_unify_post 必须为 true，可选值为 english
const UNIFY_POST_MODEL_NAME string = "unify_post_model_name"

//response parameter key
const KEY_TASK string = "Task"
const KEY_TASK_ID string = "TaskId"
const KEY_STATUS_TEXT string = "StatusText"
const KEY_RESULT string = "Result"

//status code
const STATUS_SUCCESS string = "SUCCESS"
const STATUS_RUNNING string = "RUNNING"
const STATUS_QUEUEING string = "QUEUEING"

//Initiate recording file identification
//Documentation of this interface: https://help.aliyun.com/document_detail/90727.html?spm=a2c4g.11186623.6.581.691af6ebYsUkd1
func (c AliyunClound) NewAudioFile(fileLink string) (string, *sdk.Client, error) {
	regionId, domain, apiVersion, product := c.GetApiVariable()

	fmt.Println(regionId, domain, apiVersion, product, c.Region)

	client, err := sdk.NewClientWithAccessKey(regionId, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return "", client, err
	}
	client.SetConnectTimeout(time.Second * 20)

	postRequest := requests.NewCommonRequest()
	postRequest.Domain = domain
	postRequest.Version = apiVersion
	postRequest.Product = product
	postRequest.ApiName = POST_REQUEST_ACTION
	postRequest.Method = "POST"

	mapTask := make(map[string]string)
	mapTask[KEY_APP_KEY] = c.AppKey
	mapTask[KEY_FILE_LINK] = fileLink
	// Please use version 4.0 for new access, already connected (default 2.0) if you want to maintain the status quo, please comment out this parameter setting
	mapTask[KEY_VERSION] = "4.0"
	// Set whether to output word information, the default is false, you need to set the version to 4.0 when opening
	mapTask[KEY_ENABLE_WORDS] = "true"

	//Enable unified postprocessing
	//mapTask[ENABLE_UNIFY_POST] = "true"
	//mapTask[ENABLE_INVERSE_TEXT_NORMALIZATION] = "true"
	//mapTask[UNIFY_POST_MODEL_NAME] = "chinese"

	//to json
	task, err := json.Marshal(mapTask)
	if err != nil {
		return "", client, errors.New("to json error .")
	}
	postRequest.FormParams[KEY_TASK] = string(task)
	//Make a request
	postResponse, err := client.ProcessCommonRequest(postRequest)
	if err != nil {
		return "", client, err
	}
	postResponseContent := postResponse.GetHttpContentString()
	//Check request
	if postResponse.GetHttpStatus() != 200 {
		return "", client, errors.New("Recording file identification request failed , Http error : " + strconv.Itoa(postResponse.GetHttpStatus()))
	}
	//Analytical data
	var postMapResult map[string]interface{}
	err = json.Unmarshal([]byte(postResponseContent), &postMapResult)
	if err != nil {
		return "", client, errors.New("to map struct error .")
	}

	var taskId = ""
	var statusText = ""
	statusText = postMapResult[KEY_STATUS_TEXT].(string)

	//Test result
	if statusText == STATUS_SUCCESS {
		taskId = postMapResult[KEY_TASK_ID].(string)
		return taskId, client, nil
	}

	return "", client, errors.New("Fail to detect audio files , (" + c.GetErrorStatusTextMessage(statusText) + ")")
}

//Obtain the recognition result of the recording file
//Documentation for this interface: https://help.aliyun.com/document_detail/90727.html?spm=a2c4g.11186623.6.581.691af6ebYsUkd1
func (c AliyunClound) GetAudioFileResult(taskId string, client *sdk.Client, callback func(result []byte)) error {
	_, domain, apiVersion, product := c.GetApiVariable()

	getRequest := requests.NewCommonRequest()
	getRequest.Domain = domain
	getRequest.Version = apiVersion
	getRequest.Product = product
	getRequest.ApiName = GET_REQUEST_ACTION
	getRequest.Method = "GET"
	getRequest.QueryParams[KEY_TASK_ID] = taskId
	statusText := ""

	//Traverse to get the recognition result
	for true {
		getResponse, err := client.ProcessCommonRequest(getRequest)
		if err != nil {
			return err
		}
		getResponseContent := getResponse.GetHttpContentString()

		if getResponse.GetHttpStatus() != 200 {
			return errors.New("Identification result query request failed, Http error code : " + strconv.Itoa(getResponse.GetHttpStatus()))
		}

		var getMapResult map[string]interface{}
		err = json.Unmarshal([]byte(getResponseContent), &getMapResult)
		if err != nil {
			return err
		}

		//Call callback function
		callback(getResponse.GetHttpContentBytes())

		//Check traversal conditions
		statusText = getMapResult[KEY_STATUS_TEXT].(string)
		if statusText == STATUS_RUNNING || statusText == STATUS_QUEUEING {
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	if statusText != STATUS_SUCCESS {
		return errors.New("Recording file recognition failed, (" + c.GetErrorStatusTextMessage(statusText) + ")")
	}

	return nil
}

//Get API constants
func (c AliyunClound) GetApiVariable() (regionId string, domain string, apiVersion string, product string) {
	if c.Region == 0 || c.Region == ALIYUN_CLOUND_REGION_CHA {
		regionId = REGION_ID_REGION_CHA
		domain = DOMAIN_REGION_CHA
		apiVersion = API_VERSION_REGION_CHA
		product = PRODUCT_REGION_CHA
	} else if c.Region == ALIYUN_CLOUND_REGION_INT {
		regionId = REGION_ID_REGION_INT
		domain = DOMAIN_REGION_INT
		apiVersion = API_VERSION_REGION_INT
		product = PRODUCT_REGION_INT
	}
	return
}

//get error message
func (c AliyunClound) GetErrorStatusTextMessage(statusText string) string {
	var code map[string]string = map[string]string{
		"REQUEST_APPKEY_UNREGISTERED":    "The Alibaba Cloud Smart Voice project has not been created/has no access rights. Please check whether the voice engine Appkey is incorrectly filled in; if it is an overseas region, when the software creates a voice engine, the service area needs to select "overseas"",
		"USER_BIZDURATION_QUOTA_EXCEED":  "2 hours a day to identify the free quota exceeding the limit",
		"FILE_DOWNLOAD_FAILED":           "The file access fails. Please check the OSS storage space access permission. Please set the OSS storage space to "Public Read"",
		"FILE_TOO_LARGE":                 "Audio file exceeds 512MB",
		"FILE_PARSE_FAILED":              "Audio file parsing failed, please check the audio file for corruption",
		"UNSUPPORTED_SAMPLE_RATE":        "Sample rate mismatch",
		"FILE_TRANS_TASK_EXPIRED":        "Audio file recognition task expired, please try again",
		"REQUEST_INVALID_FILE_URL_VALUE": "Audio file access failed, please check the OSS storage space access permission",
		"FILE_404_NOT_FOUND":             "Audio file access failed, please check the OSS storage space access permission",
		"FILE_403_FORBIDDEN":             "Audio file access failed, please check the OSS storage space access permission",
		"FILE_SERVER_ERROR":              "Audio file access failed, please check if the service where the requested file is located is available",
		"INTERNAL_ERROR":                 "Internal generic error identified, please try again later",
	}

	if _, ok := code[statusText]; ok {
		return code[statusText]
	} else {
		return statusText
	}
}
