package main

import (
	"github.com/gin-gonic/gin"
	"github.com/smtc/glog"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func assetsFiles(c *gin.Context) {
	r := c.Request
	pth := c.Param("pth")
	if pth == "" {
		glog.Error("assetsFiles: path is empty: %s\n", r.URL.Path)
		c.Data(200, "text/plain", []byte(""))
		return
	}

	fp, err := getAssetFilePath(pth)
	if err != nil {
		glog.Error("assetsFiles: %s\n", err)
		c.Data(200, "text/plain", []byte(""))
		return
	}

	http.ServeFile(c.Writer, c.Request, fp)
}

func getAssetFilePath(pth string) (string, error) {
	entrys := strings.Split(pth, "/")
	sentrys := []string{contentDir}
	for _, s := range entrys {
		s = strings.TrimSpace(s)
		if s != "" {
			sentrys = append(sentrys, s)
		}
	}
	return path.Join(sentrys...), nil
}

/**
图片文件提取
创建人:邵炜
创建时间:2016年3月7日11:47:17
输入参数: gin对象
输出参数: 无
数据反馈由gin对象进行
*/
func imageGet(c *gin.Context) {
	idStr := c.Param("name")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		glog.Error("imageGet param is error, err: %s \n", err.Error())
		c.String(http.StatusOK, "image bytes is empty")
		return
	}

	_imageFile.Lock()

	fileName,ok:=_imageFile.imageFile[id]

	defer _imageFile.Unlock()

	if !ok {
		glog.Error("imageGet _imageFile array can't find file , id: %s \n",id)
		return
	}

	http.ServeFile(c.Writer,c.Request,fileName.pictureName)
}
