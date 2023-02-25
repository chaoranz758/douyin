package gin_middleware

//func Ja3() gin.HandlerFunc {
//	return func(context *gin.Context) {
//		ja3ContextData := ja3.GetRequestJa3Data(context.Request)
//		//返回设备名称，判断是否是安全流量
//		_, ja3OK := ja3ContextData.Verify()
//		if !ja3OK {
//			context.JSON(http.StatusOK, gin.H{
//				"status_code": handler.CodeCrawlerRequest,
//				"msg":         "不安全的访问,请使用浏览器访问,这次请求将被拦截",
//			})
//			context.Abort()
//			return
//		}
//		context.Next()
//	}
//}
