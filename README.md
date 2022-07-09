# lectureNoteAssistant
A windows desktop application which can generate subtitles and translations for videos which you can use for generating billingual transcripts for videos.

## How to use it?

Just download the [exe](https://github.com/WideSu/lectureNoteAssistant/blob/main/lectureNoteAssistant.exe) file, and run it on your **windows** computer.


## Supported video types
Supports video and audio files in common formats, including: 
- For videos, we support: .mp4 , .mpeg , .mkv , .wmv , .avi , .m4v , .mov , .flv , .rmvb , .3gp , .f4v . 
- For audio: .mp3 , .wav , .aac , .wma , .flac , .m4a formats.

## Output files includes
It can generate 3 types Support subtitle files including: SRT file, ordinary text, lecture summary simultaneously or seperately. And with bilingual translation between 10 languages including Chinese, English, Japanese, Korean, French, German, Spanish, Russian, Italian, and Thai. 

## Services and SDKs used:
- Baidu and tencent Translation SDK
- go tldr for auto-summary
- **aliyun-cloud-sdk-go，aliyun-oss-go-sdk，tencentcloud-sdk-go** for storing the audio files and generated transcripts
- the intelligent voice interactive service on Aliyun

## Demo of this app

- Firstly, upload the video or audio on your computer

<img src="https://github.com/WideSu/lectureNoteAssistant/blob/main/screenshot/LectureNoteAssistant_1.gif">

- Secondly, submit it. The software will start to generate transcripts for you

<img src="https://github.com/WideSu/lectureNoteAssistant/blob/main/screenshot/lectureNoteAssistant.gif">
