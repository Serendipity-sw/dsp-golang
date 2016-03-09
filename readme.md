#web框架

## 框架须知请参考

框架项目目标地址: [亲,touch me](https://github.com/swgloomy/webframe)

## 项目目前终止更新

    因业务扩展原因,该项目目前终止更新,框架项目会不定时进行更新
    
## 项目说明描述
    
    1.增加图片数据库读取,并存放到临时目录及临时目录删除方法,方法定时执行,无需手动更改  方法存放在 timeMonitor.go文件
    
    2.增加路由 /assets/*path  检索资源文件,方法为assetsFiles 寄存在router.go文件
    
    3.增加路由 /images/:name  图片资源文件获取  方法为 imageGet  寄存在router.go文件