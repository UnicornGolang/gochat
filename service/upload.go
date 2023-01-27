package service

import (
	"fmt"
	"gochat/utils"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
  writer := c.Writer
  request := c.Request
  fileType := c.PostForm("filetype")
  srcFile, head, err := request.FormFile("file")

  if err != nil {
    utils.RespFail(writer, err.Error())
    return
  }
  suffix := ".png"
  ofileName := head.Filename
  fmt.Println("上传文件名", ofileName)
  if fileType != "" {
    suffix = fileType
  }else {
    tmp := strings.Split(ofileName, ".")
    if len(tmp) > 1 {
      suffix = "." + tmp[len(tmp)-1]
    }
  }
  fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
  filePath := "./asset/upload/" + fileName
  dstFile, err := os.Create(filePath)
  if err != nil {
    utils.RespFail(writer, err.Error())
    return
  }
  _, err = io.Copy(dstFile, srcFile)
  if err != nil {
    utils.RespFail(writer, err.Error())
    return
  }
  // 将文件的存储路径返回出来
  utils.RespOk(writer, filePath, "上传成功")
}
