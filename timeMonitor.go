package main

import (
	"github.com/gin-gonic/gin"
	"github.com/guotie/deferinit"
	"github.com/howeyc/fsnotify"
	"github.com/smtc/glog"
	"html/template"
	"sync"
	"time"
	"github.com/guotie/config"
"strings"
	"encoding/base64"
	"fmt"
	"os"
	"io"
	"io/ioutil"
)

type imageFileStruct struct {
		pictureName string
}

type imageFileMap struct {
	sync.RWMutex
	imageFile map[int]imageFileStruct
}

var (
	jsTmr    *time.Timer
	funcName = template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeurl": func(s string) template.URL {
			return template.URL(s)
		},
	}

	_imageFile=imageFileMap{
		imageFile:map[int]imageFileStruct{},
	}
)

func init() {
	deferinit.AddInit(func() {
		tempDir = config.GetStringDefault("tempDir", "template/")
		if !strings.HasSuffix(tempDir, "/"){
			tempDir += "/"
		}
	}, nil, 40)
	deferinit.AddRoutine(notifyTemplates)
	deferinit.AddRoutine(watchJsDir)
}

/**
定时运行程序
创建人:邵炜
创建时间:2016年3月7日09:51:42
输入参数: 终止命令  计数器对象
*/
func watchJsDir(ch chan struct{}, wg *sync.WaitGroup) {
	go func() {
		<-ch

		jsTmr.Stop()
		wg.Done()
	}()

	jsTmr = time.NewTimer(time.Minute)
	for {
		imageFile()
		jsTmr.Reset(time.Minute)
		<-jsTmr.C
	}
}

/**
加载模版
创建人:邵炜
创建时间:2016年2月26日11:34:12
输入参数: gin对象
*/
func loadTemplates(e *gin.Engine) {
	t, err := template.New("tmpls").Funcs(funcName).ParseGlob(tempDir+"*")

	if err != nil {
		glog.Error("loadTemplates failed: %s %s \n", tempDir, err.Error())
		return
	}

	e.SetHTMLTemplate(t)
}

/**
监视文件夹目录如发生任何修改,重新载入
创建人:邵炜
创建时间:2016年3月7日09:47:50
输入参数: 终止命令 计数器对象
*/
func notifyTemplates(ch chan struct{}, wg *sync.WaitGroup) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Error("notifyTemplates: create new watcher failed: %v\n", err)
		return
	}

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				glog.Debug("notifyTemplates: event: %v\n", ev)
				loadTemplates(rt)
			case err := <-watcher.Error:
				glog.Error("notifyTemplates: error: %v\n", err)
			}
		}
	}()

	err = watcher.Watch(tempDir)
	if err != nil {
		glog.Error("notifyTemplates: watch dir %s failed: %v \n", tempDir, err)
	}

	// Hang so program doesn't exit
	<-ch

	/* ... do clean stuff ... */
	watcher.Close()
	wg.Done()
}

/**
图片流处理
创建人:邵炜
创建时间:2016年3月8日14:20:41
 */
func imageFile() {
	imageArray,err:=selectImageAll()

	if err != nil {
		glog.Error("imageFile selectImage is error, err: %s \n",err.Error())
		return
	}

	_imageFile.Lock()

	_imageFile.imageFile=map[int]imageFileStruct{}

	imageFileDelete()

	defer _imageFile.Unlock()

	for _,value:=range *imageArray  {

		images:=strings.Split(value.pictureByte,",")

		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(images[1]))

		imageName:=strings.Split(strings.Split(images[0],";")[0],"/")[1]

		fileName:=fmt.Sprintf(temp+"/%d%d.%s",time.Now().Unix(),value.id,imageName)

		file, err := os.Create(fileName)

		if err != nil {
			glog.Error("image create is error, err: %s \n",err.Error())
			file.Close()
			continue
		}

		_, err = io.Copy(file,reader)

		if err != nil {
			glog.Error("image copy is error, err: %s \n",err.Error())
			file.Close()
			continue
		}

		_imageFile.imageFile[value.id]=imageFileStruct{
			pictureName:fileName,
		}

		file.Close()
	}

	glog.Info("image save computer file is success \n")
}

/**
图片资源清除
创建人:邵炜
创建时间:2016年3月8日16:43:10
 */
func imageFileDelete() {
	files,err:=ioutil.ReadDir(temp+"/")

	if err != nil {
		glog.Error("imageFileDelete can't search files, error: %s \n",err.Error() )
		return
	}

	for _,file:=range files  {
		if file.IsDir() {
			continue
		}

		fileName:=fmt.Sprintf("%s/%s",temp,file.Name())

		err=os.Remove(fileName)

		if err != nil {
			glog.Error("imageFileDelete can't delete file, fileName: %s  err: %s \n",fileName,err.Error())
		}
	}

	glog.Info("imageFileDelete delete file is perfection! \n")
}
