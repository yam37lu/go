package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.b.Write(b)
	return w.ResponseWriter.Write(b)
}

//func OperateLog() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		start := time.Now()
//		writer := responseWriter{
//			c.Writer,
//			bytes.NewBuffer([]byte{}),
//		}
//		c.Writer = writer
//		c.Next()
//		tmp := strings.Split(c.Request.URL.Path, "/")
//		operateType, ok := global.SYS_CONFIG.System.OperateTypeMap[strings.Join(tmp[3:], "/")]
//		if ok {
//			resp := response.Response{}
//			if err := json.Unmarshal(writer.b.Bytes(), &resp); err == nil {
//				state := 1
//				if resp.Code != utils.Normal {
//					state = 0
//				}
//				dataIds := strings.Split(c.GetString("dataId"), ",")
//				logs := make([]*model.OperateLog, 0)
//				for _, id := range dataIds {
//					td := strings.Split(id, ":")
//					if len(td) == 2 {
//						t1, err := strconv.Atoi(td[0])
//						if err != nil {
//							global.SYS_LOG.Error(fmt.Sprintf("operate log, dataId[%s] error[%v]", id, err))
//							continue
//						}
//						operateType = t1
//						id = td[1]
//					}
//					logs = append(logs, &model.OperateLog{
//						Id:         utils.UUID(),
//						DataId:     id,
//						State:      state,
//						Type:       operateType,
//						Message:    resp.Message,
//						Result:     writer.b.String(),
//						Operator:   c.GetString("user"),
//						CreateTime: start,
//						ModifyTime: time.Now(),
//					})
//				}
//				dao.AddLogs(logs)
//			}
//		}
//	}
//}
