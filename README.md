# lectureNoteAssistant

A Windows desktop application that can generate subtitles, translations, and summaries for videos in 8 languages using API and SDK from Tencent, Alibaba, and Baidu. You can use it for generating bilingual transcripts for videos and summarising the key points from the transcript using the PageRank algorithm: LexRank.

<img width="895" alt="image" src="https://user-images.githubusercontent.com/44923423/178111967-1675c55f-501e-4f60-8d07-fe6fc023c60b.png">


## How to use it?

Just download the [exe](https://github.com/WideSu/lectureNoteAssistant/blob/main/lectureNoteAssistant.exe) file, and run it on your **windows** computer.


## Supported video types
Supports video and audio files in common formats, including: 
- For videos, we support: .mp4 , .mpeg , .mkv , .wmv , .avi , .m4v , .mov , .flv , .rmvb , .3gp , .f4v . 
- For audio: .mp3 , .wav , .aac , .wma , .flac , .m4a formats.

## Output files
It can generate 3 types Support subtitle files including: SRT file, ordinary text, lecture summary simultaneously or seperately according to user's needs. And with bilingual translation between 10 languages including Chinese, English, Japanese, Korean, French, German, Spanish, Russian, Italian, and Thai. 

## Services and SDKs used:
- Baidu and tencent Translation SDK
- go tldr for auto-summary
- **aliyun-cloud-sdk-go，aliyun-oss-go-sdk，tencentcloud-sdk-go** for storing the audio files and generated transcripts
- the intelligent voice interactive service on Aliyun

## Demo

- Firstly, upload the video or audio on your computer

<img src="https://github.com/WideSu/lectureNoteAssistant/blob/main/screenshot/LectureNoteAssistant_1.gif">

- Secondly, submit it. The software will start to generate transcripts for you

<img src="https://github.com/WideSu/lectureNoteAssistant/blob/main/screenshot/lectureNoteAssistant.gif">

## System architecture

This system uses the go [walk](https://github.com/lxn/walk) library for development, and main.go contains the code and main logic of the main interface of the program. It calls the relevant code files in the app package to perform corresponding operations. Interface logic code, data object separation. Basically similar to the MVC pattern. It is divided into presentation layer, business logic layer and data access layer. The presentation layer is used to interact with the user, and then calls the functions of the corresponding modules in the app package to perform business operations, and the corresponding modules of the business operations then call the data layer functions to operate on the data.

<img width="1086" alt="image" src="https://user-images.githubusercontent.com/44923423/178111178-3838b4de-4663-4d2e-b79d-b8fc01a20885.png">


