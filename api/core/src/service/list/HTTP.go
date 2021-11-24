package member

import (
	"context"
	"core/src/code"
	"core/src/db/mariadb"
	"core/src/lib/errs"
	"core/src/module/log"
	"github.com/gin-gonic/gin"
	"time"
)

type HTTP struct {
}

func (h *HTTP) Mount(ctx context.Context, engine *gin.Engine) error {
	router := engine.Group("/project_cards")
	router.GET("/project_summary", h.projectSummary(ctx))
	return nil
}

func (h *HTTP) projectSummary(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		type request struct {
			StartTime int64  `json:"start_time"`
		}
		var req request
		err := c.ShouldBindJSON(&req)
		if err != nil {
			code.MissingInputValue.JsonResponse(c, err)
			return
		}

		cardList, err := getProjectSummary(ctx, req.StartTime)
		if err != nil {
			code.ServerError.JsonResponse(c, err)
			return
		}

		code.OK.JsonResponse(c, cardList)
	}
}

type ProjectCard struct {
	UUID       string    `json:"uuid"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	ContentUrl string    `json:"content_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

//
//func init() {
//	err := syncProjectSummary()
//	if err != nil{
//		panic(err)
//	}
//}
//
//func syncProjectSummary()error {
//
//	cardList, err := getProjectSummaryInDB(conf.GlobalConfig.)
//	if err != nil {
//		if e1 != nil {
//			err = errs.Join(e1, err)
//		}
//		err = errs.WithLine(err)
//		return
//	}
//}

func getProjectSummary(ctx context.Context, updateTime int64) (cardList []ProjectCard, err error) {
	cardList, err = getProjectSummaryInDB(ctx, updateTime)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	return cardList, err

	//var ok bool
	//var e1 error
	//
	//cardList, ok, e1 = getProjectSummaryInCache(ctx,updateTime)
	//if e1 != nil {
	//	e1 = errs.WithLine(e1)
	//}
	//if ok {
	//	return
	//}

	//cardList, err = getProjectSummaryInDB(ctx,updateTime)
	//if err != nil {
	//	if e1 != nil {
	//		err = errs.Join(e1, err)
	//	}
	//	err = errs.WithLine(err)
	//	return
	//}
	//err = cacheProjectSummary(ctx,cardList)
	return

}

//const summaryContentKey = "summary_content_key"

//func cacheProjectSummary(ctx context.Context,list []ProjectCard) error {
//	pp := redisdb.Conn.Pipeline()
//	for _, card := range list {
//		_, err := pp.ZAdd(ctx,card.Title , &redis.Z{
//			Score:  float64( card.UpdatedAt.Unix()),
//			Member: card.UUID,
//		}).Result()
//		if err !=  nil{
//			err = errs.WithLine(err)
//			return err
//		}
//		pp.HMSet(summaryContentKey , )
//	}
//
//
//	//fmt.Println("updateTime: ",updateTime)
//	//
//	//b := make([]byte, 8)
//	//binary.LittleEndian.PutUint64(b, uint64(updateTime))
//	//
//	//fmt.Println("b: ",b)
//	//
//	//
//	//i := int64(binary.LittleEndian.Uint64(b))
//	//fmt.Println("i: ",i)
//
//	//_, err := redisdb.Conn.HSetStruct(ctx , redisdb.FormatKey(projectSummaryKey) ,updateTime,list)
//	//if err != nil {
//	//	return err
//	//}
//	//return nil
//}
//
//func getProjectSummaryInCache(ctx context.Context,updateTime int64) ([]ProjectCard, bool, error) {
//	var pc  []ProjectCard
//	err := redisdb.Conn.HScan()
//	if err != nil {
//		if err == redis.Nil{
//			return nil, false, nil
//		}
//		return nil, false, err
//	}
//	return pc , true , nil
//}

func getProjectSummaryInDB(ctx context.Context, updateTime int64) ([]ProjectCard, error) {
	res, err := mariadb.Conn.QueryContext(ctx, "select title , summary , content_url , updated_at,created_at from project_card_list where updated_at > FROM_UNIXTIME(?);", updateTime)
	if err != nil {
		err = errs.WithLine(err)
		return nil, err
	}
	defer errs.LogIfErr(res.Close, log.Logger.Skip(1).Error)

	var retList []ProjectCard

	for res.Next() {
		row := ProjectCard{}
		e := res.Scan(&row.Title, &row.Summary, &row.ContentUrl, row.UpdatedAt, &row.CreatedAt)
		if e != nil {
			e = errs.WithLine(e)
			return nil, e
		}
		retList = append(retList, row)
	}
	return retList, nil
}
