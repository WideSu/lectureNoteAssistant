# LectureNoteAssistant

## 简介

`LectureNoteAssistant` 是采用`Golang`语言开发的一个Windows GUI桌面应用。<br /> 可以用于识别中英文讲座视频语音自动生成字幕SRT/歌词LRC/文本TXT文件，以及视频摘要。<br />

本应用使用了以下接口：
- 阿里云 [OSS对象存储](https://www.aliyun.com/product/oss?spm=5176.12825654.eofdhaal5.13.e9392c4aGfj5vj&aly_as=K11FcpO8)
- 阿里云 [录音文件识别](https://ai.aliyun.com/nls/filetrans?spm=5176.12061031.1228726.1.47fe3cb43I34mn) 
- 百度翻译开放平台 [翻译API](http://api.fanyi.baidu.com/api/trans/product/index) 
- 腾讯云 [翻译API](https://cloud.tencent.com/product/tmt) 

<a name="0b884e4f"></a>
## 软件界面

![image](https://github.com/WideSu/lectureNoteAssistant/blob/main/screenshot/lectureNoteAssistant.gif)

## 应用场景

- 识别**中英文教育视频/音频**的语音生成字幕文件（支持中英互译，双语字幕）
- 提取**视频/音频**的语音文本
- 批量翻译、过滤处理/编码SRT字幕文件
- 自动生成视频摘要，方便管理讲座录屏


<a name="b89d37d3"></a>
## 软件优势

- 使用阿里云智能语音识别接口，准确度高，标准普通话/英语识别率95%以上
- 视频识别无需上传原视频，只需上传音频文件至阿里云接口，方便快速且节省时间
- 使用阿里云OSS对象储存生成的字幕文件，无需本地数据库
- 支持多任务多文件批量处理
- 支持视频、音频常见多种格式文件
（支持的视频格式：.mp4 , .mpeg , .mkv , .wmv , .avi , .m4v , .mov , .flv , .rmvb , .3gp , .f4v
  支持的音频格式：.mp3 , .wav , .aac , .wma , .flac , .m4a
  支持的字幕格式：.srt）
- 支持同时输出字幕SRT文件、LRC文件、普通文本3种类型
- 支持语气词过滤、自定义文本过滤、正则过滤等，使软件生成的字幕更加精准
- 支持字幕中英互译、双语字幕输出，及日语、韩语、法语、德语、西班牙语、俄语、意大利语、泰语等
- 支持多翻译引擎（百度翻译、腾讯云翻译）
- 支持批量翻译、编码SRT字幕文件

<a name="1bbbb204"></a>
## 注意事项

- 软件目录下的 `data`目录为数据存储目录，请勿删除。否则可能会导致配置丢失
- 项目使用了 [ffmpeg](http://ffmpeg.org/) 依赖，如果电脑没有配置过ffmpeg，需要在ffmpeg官网下了一个full-build的版本，解压后将其中bin目录加入系统变量Path中

## FAQ

##### 1.使用此软件会产生费用吗？
如果您适量使用本软件（各个API的免费使用额度可以自行查询，如阿里云语音识别免费版每天限量2h），将不会产生费用。
如果您大量使用，建议根据自己的情况购买各个平台的资源包，以满足需求。

##### 2.为何软件一直报错？
报错的原因有很多，未配置ffmpeg依赖、软件运行命令错误、阿里云、腾讯云（注意要使用共网链接）等账户权限问题都可能会导致软件显示错误。如果您遇到麻烦，可以加QQ 1197749338与我交流。

##### 3.如何运行？
1. 在Go官网下载安装包，参考[Go安装教程](https://golang.org/doc/install)<br />
2. 在ffmpeg官网下载对应操作系统的full-build的安装包，解压后将其中bin文件加入系统变量Path中[ffmpeg下载地址](https://ffmpeg.org/download.html)<br />
3. 在VS Code中配置Go开发环境（Go扩展 以及launch.json）[VS code配置Go开发环境教程](https://www.liwenzhou.com/posts/Go/00_go_in_vscode/)<br />
4. 导入本项目<br />
5. 在终端输入go build -ldflags="-H windowsgui"编译项目产生可执行文件，参考[walk官方教程](https://github.com/lxn/walk)<br />
6. 运行可执行文件<br />
7. 配置阿里云语音接口[录音文件识别](https://ai.aliyun.com/nls/filetrans?spm=5176.12061031.1228726.1.47fe3cb43I34mn) <br />
8. 配置腾讯、百度翻译接口[百度翻译API](http://api.fanyi.baidu.com/api/trans/product/index)[腾讯翻译API](https://cloud.tencent.com/product/tmt)  <br />
9. 配置阿里云OSS存储接口[OSS对象存储](https://www.aliyun.com/product/oss?spm=5176.12825654.eofdhaal5.13.e9392c4aGfj5vj&aly_as=K11FcpO8)<br />

##### 4.如何提升语音识别准确度？
1. 可在阿里云智能语音交互平台上上传对应的训练集，训练自定义模型[阿里云智能语音交互自定义模型官方Doc](https://help.aliyun.com/document_detail/72216.html?spm=a2c4g.11186623.6.565.3d0569386dk3T3)<br />
2. 可在阿里云智能语音交互平台上添加专业词汇热词，改善软件对领域专业词汇的识别准确度[阿里云智能语音交互管理热词官方Doc](https://help.aliyun.com/document_detail/72215.html?spm=a2c4g.11186623.6.564.40071037R34ic5)<br />
3. 在软件中设置语气词、替换规则<br />

<a name="f3dc992e"></a>

## 测试视频来源：LED演讲视频
链接：https://pan.baidu.com/s/16gfn1lUSNiNWKCEP1uUlKQ <br />
提取码：fjq4 <br />
复制这段内容后打开百度网盘手机App，操作更方便哦--来自百度网盘超级会员V3的分享<br />

## 交流&联系

- QQ：1197749338

