package main

import (
	"github.com/smtc/glog"
)

type imageAllStruct struct {
	pictureByte string
	id int
}

/**
根据ID查询图片
创建人:邵炜
创建时间:2016年3月7日11:53:59
输入参数: 图片主键
输出参数: 图片二进制字符串 错误对象
*/
func selectImageById(id int) (*string, error) {
	sqlStr := "select PictureByte from pictures where Id=?;"

	rows, err := sqlSelect(sqlStr, id)

	if err != nil {
		glog.Error("selectImage is error, sqlStr: %s id: %d err: %s \n", sqlStr, id, err.Error())
		return nil, err
	}

	defer rows.Close()

	var imageByte string

	for rows.Next() {
		err = rows.Scan(&imageByte)
		if err != nil {
			glog.Error("selectImage PictureByte is error , sqlStr: %s id: %d err: %s \n", sqlStr, id, err.Error())
			return nil, err
		}
		break
	}

	glog.Info("selectImage is success! \n")

	return &imageByte, nil
}

/**
查询所有的图片
创建人:邵炜
创建时间:2016年3月8日15:25:57
输入参数: 无
输出参数: 图片数组对象 错误对象
 */
func selectImageAll() (*[]imageAllStruct, error) {
	sqlStr:="select Id,PictureByte from pictures"

	rows,err:=sqlSelect(sqlStr)

	if err != nil {
		glog.Error("selectImageAll is error, sqlStr: %s err: %s \n",sqlStr,err.Error())
		return nil,err
	}

	var imageArray []imageAllStruct

	for rows.Next()  {
		var(
			id int
			imageByte string
		)

		err=rows.Scan(&id,&imageByte)

		if err != nil {
			glog.Error("selectImageAll read error,sqlStr: %s err: %s \n",sqlStr,err.Error())
			continue
		}

		imageArray=append(imageArray,imageAllStruct{
			id:id,
			pictureByte:imageByte,
		})

	}

	glog.Info("selectImageAll is success! \n")

	return &imageArray,nil
}
