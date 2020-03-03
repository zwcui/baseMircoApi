/*
@Time : 2019/12/16 下午12:00 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"encoding/json"
	"jingting_server/publicservice/models"
	"jingting_server/publicservice/util"
	"jingting_server/publicservice/base"
	"strconv"
	"jingting_server/publicservice/remote"
)

type QuestionnaireController struct {
	apiController
}

func (this *QuestionnaireController) Prepare(){
	this.NeedAuthorAuthList = []RequestPathAndMethod{

	}
	this.authorAuth()
	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/getQuestionnaireAuthorStatisticsList", "get", []int{1,2,3,4}},
		{"/getQuestionnaireAnswerList", "get", []int{1,2,3,4}},
	}
	this.userAuth()
}

// @Title 上传回答解析（h5使用）
// @Description 上传回答解析（h5使用）
// @Param	questionnaireId			formData		int64	  		true		"questionnaireId"
// @Param	answerJson				formData		string	  		true		"answerJson"
// @Param	phoneNumber				formData		string	  		false		"phoneNumber"
// @Success 200 {string} success
// @router /answerQuestionnaire [post]
func (this *QuestionnaireController) AnswerQuestionnaire() {
	questionnaireId := this.MustInt64("questionnaireId")
	answerJson := this.MustString("answerJson")
	phoneNumber := this.GetString("phoneNumber")

	util.Logger.Info(answerJson)

	var answerList []models.QuestionnaireAnswerRequest
	err := json.Unmarshal([]byte(answerJson), &answerList)
	if err != nil {
		util.Logger.Info("err:" + err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.QuestionnaireError100)
		return
	}

	if answerList != nil && len(answerList) > 0 {
		if answerList[0].AuthorId == 0 {
			this.ReturnData = util.GenerateAlertMessage(models.QuestionnaireError600)
			return
		}
	}

	directionMap := make(map[string][]int)
	directionWeightMap := make(map[string]float64)

	var questionnaireAuthorStatistics models.QuestionnaireAuthorStatistics
	base.DBEngine.Table("questionnaire_author_statistics").InsertOne(&questionnaireAuthorStatistics)
	for _, answer := range answerList {
		if answer.QuestionnaireQuestionId == 0 {
			switch answer.Direction {
			case "sex":
				questionnaireAuthorStatistics.Sex = answer.Answer
				break
			case "age":
				questionnaireAuthorStatistics.Age = answer.Answer
				break
			case "education":
				questionnaireAuthorStatistics.Education = answer.Answer
				break
			case "city":
				questionnaireAuthorStatistics.City = answer.Answer
				break
			case "occupation":
				questionnaireAuthorStatistics.Occupation = answer.Answer
				break
			}
			continue
		}

		var questionnaireAnswer models.QuestionnaireAnswer
		questionnaireAnswer.QuestionnaireQuestionId = answer.QuestionnaireQuestionId
		questionnaireAnswer.QuestionnaireAuthorStatisticsId = questionnaireAuthorStatistics.Id
		questionnaireAnswer.AuthorId = answer.AuthorId
		questionnaireAnswer.Direction = answer.Direction
		questionnaireAnswer.Content = answer.Content
		questionnaireAnswer.Answer = answer.Answer
		questionnaireAnswer.Score = answer.Score
		base.DBEngine.Table("questionnaire_answer").InsertOne(&questionnaireAnswer)

		_, ok := directionMap[answer.Direction]
		if ok {
			directionMap[answer.Direction] = append(directionMap[answer.Direction], answer.Score)
		} else {
			var scoreArray []int
			scoreArray = append(scoreArray, answer.Score)
			directionMap[answer.Direction] = []int(scoreArray)
		}

		var questionnaireQuestion models.QuestionnaireQuestion
		base.DBEngine.Table("questionnaire_question").Where("questionnaire_id=?", questionnaireId).And("direction=?", answer.Direction).Get(&questionnaireQuestion)
		directionWeightMap[answer.Direction] = questionnaireQuestion.DirectionWeight
	}

	util.Logger.Info("directionMap")
	util.Logger.Info(directionMap)

	index := 1
	questionnaireAuthorStatistics.PhoneNumber = phoneNumber
	for direction, scoreArray := range directionMap {
		questionnaireAuthorStatistics.AuthorId = answerList[0].AuthorId
		questionnaireAuthorStatistics.QuestionnaireId = questionnaireId
		switch index {
		case 1:
			questionnaireAuthorStatistics.Direction1 = direction
			questionnaireAuthorStatistics.Average1 = calcScoreAverage(scoreArray)
			break
		case 2:
			questionnaireAuthorStatistics.Direction2 = direction
			questionnaireAuthorStatistics.Average2 = calcScoreAverage(scoreArray)
			break
		case 3:
			questionnaireAuthorStatistics.Direction3 = direction
			questionnaireAuthorStatistics.Average3 = calcScoreAverage(scoreArray)
			break
		case 4:
			questionnaireAuthorStatistics.Direction4 = direction
			questionnaireAuthorStatistics.Average4 = calcScoreAverage(scoreArray)
			break
		case 5:
			questionnaireAuthorStatistics.Direction5 = direction
			questionnaireAuthorStatistics.Average5 = calcScoreAverage(scoreArray)
			break
		case 6:
			questionnaireAuthorStatistics.Direction6 = direction
			questionnaireAuthorStatistics.Average6 = calcScoreAverage(scoreArray)
			break
		case 7:
			questionnaireAuthorStatistics.Direction7 = direction
			questionnaireAuthorStatistics.Average7 = calcScoreAverage(scoreArray)
			break
		}
		index ++
	}

	questionnaireAuthorStatistics.TotalScore = util.Rounding(questionnaireAuthorStatistics.Average1 * directionWeightMap[questionnaireAuthorStatistics.Direction1] +
		questionnaireAuthorStatistics.Average2 * directionWeightMap[questionnaireAuthorStatistics.Direction2] +
		questionnaireAuthorStatistics.Average3 * directionWeightMap[questionnaireAuthorStatistics.Direction3] +
		questionnaireAuthorStatistics.Average4 * directionWeightMap[questionnaireAuthorStatistics.Direction4] +
		questionnaireAuthorStatistics.Average5 * directionWeightMap[questionnaireAuthorStatistics.Direction5] +
		questionnaireAuthorStatistics.Average6 * directionWeightMap[questionnaireAuthorStatistics.Direction6] +
		questionnaireAuthorStatistics.Average7 * directionWeightMap[questionnaireAuthorStatistics.Direction7])
	base.DBEngine.Table("questionnaire_author_statistics").Where("id=?", questionnaireAuthorStatistics.Id).AllCols().Update(&questionnaireAuthorStatistics)

	this.ReturnData = "success"
}


// @Title 查看调查问卷（h5使用）
// @Description 查看调查问卷（h5使用）
// @Param	questionnaireId			formData		int64	  		true		"questionnaireId"
// @Success 200 {object} models.QuestionnaireQuestionListContainer
// @router /showQuestionnaire [get]
func (this *QuestionnaireController) ShowQuestionnaire() {
	questionnaireId := this.MustInt64("questionnaireId")

	var questionnaireQuestionList []models.QuestionnaireQuestion
	base.DBEngine.Table("questionnaire_question").Where("questionnaire_id=?", questionnaireId).Asc("sort_no").Find(&questionnaireQuestionList)

	if questionnaireQuestionList == nil {
		questionnaireQuestionList = make([]models.QuestionnaireQuestion, 0)
	}

	this.ReturnData = models.QuestionnaireQuestionListContainer{questionnaireQuestionList}
}

// @Title 查看调查问卷结果（h5使用）
// @Description 查看调查问卷结果（h5使用）
// @Param	questionnaireId			query			int64	  		true		"questionnaireId"
// @Param	authorId				query			int64	  		true		"authorId"
// @Success 200 {object} models.QuestionnaireAuthorStatisticsContainer
// @router /getQuestionnaireResult [get]
func (this *QuestionnaireController) GetQuestionnaireResult() {
	questionnaireId := this.MustInt64("questionnaireId")
	authorId := this.MustInt64("authorId")

	var questionnaireAuthorStatistics models.QuestionnaireAuthorStatistics
	hasResult, _ := base.DBEngine.Table("questionnaire_author_statistics").Where("questionnaire_id=?", questionnaireId).And("author_id=?", authorId).Desc("created").Limit(1, 0).Get(&questionnaireAuthorStatistics)
	if !hasResult {
		this.ReturnData = util.GenerateAlertMessage(models.QuestionnaireError200)
		return
	}
	if questionnaireAuthorStatistics.Direction7 == "" {
		questionnaireAuthorStatistics.Direction7 = "家庭与子女"
	}
	var freeCourseConfig models.SystemConfig
	var course models.Course
	base.DBEngine.Table("system_config").Where("program='questionnaire_free_course_id'").Get(&freeCourseConfig)
	base.DBEngine.Table("course").Where("id = ?",freeCourseConfig.ProgramValue).Get(&course)

	this.ReturnData = models.QuestionnaireAuthorStatisticsContainer{questionnaireAuthorStatistics, course.Title, course.Id}
}

// @Title 回答问卷成功后领取课程（h5使用）
// @Description 回答问卷成功后领取课程（h5使用）
// @Param	authorId				formData		int64	  		true		"用户id"
// @Success 200 {string} success
// @router /getFreeCourseAfterQuestionnaire [post]
func (this *QuestionnaireController) GetFreeCourseAfterQuestionnaire() {
	authorId := this.MustInt64("authorId")

	var author models.Author
	hasAuthor, _ := base.DBEngine.Table("author").Where("id=?", authorId).Get(&author)
	if !hasAuthor {
		this.ReturnData = util.GenerateAlertMessage(models.QuestionnaireError300)
		return
	}

	var questionnaireAuthorStatistics models.QuestionnaireAuthorStatistics
	hasQuestionnaireAuthorStatistics, _ := base.DBEngine.Table("questionnaire_author_statistics").Where("author_id=?", authorId).Get(&questionnaireAuthorStatistics)
	if !hasQuestionnaireAuthorStatistics {
		this.ReturnData = util.GenerateAlertMessage(models.QuestionnaireError400)
		return
	}

	var freeCourseConfig models.SystemConfig
	base.DBEngine.Table("system_config").Where("program='questionnaire_free_course_id'").Get(&freeCourseConfig)

	var payOrder models.PayOrder
	hasPayOrder, _ := base.DBEngine.Table("pay_order").Where("order_type=3").And("order_type_id=?", freeCourseConfig.ProgramValue).And("payer_id=?", authorId).And("status=1").Get(&payOrder)
	if hasPayOrder {
		this.ReturnData = util.GenerateAlertMessage(models.QuestionnaireError500)
		return
	}

	payOrder.OrderType = 3
	payOrder.AppType = -2
	payOrder.OrderTypeId, _ = strconv.ParseInt(freeCourseConfig.ProgramValue, 10, 64)
	payOrder.PayerId = authorId
	payOrder.Status = 1
	base.DBEngine.Table("pay_order").InsertOne(&payOrder)
	//回答成功后推送消息
	var course models.Course
	base.DBEngine.Table("course").Where("id = ?", freeCourseConfig.ProgramValue).Get(&course)
	appMessageOne := models.AppMessage{}
	appMessageOne.ReceiverId = authorId
	appMessageOne.ActionUrl = remote.JumpUrlWithKeyAndPramas(models.JTSHARE_JUMP_KEY, nil)
	appMessageOne.Content = "感谢您认真填写问卷，赠送您一套课程《" + course.Title + "》，欢迎观看"
	appMessageOne.Type = 7
	_, pushErr := remote.PushMessageToUser(authorId, &appMessageOne, "", 0)
	if pushErr != nil {
		util.Logger.Info("addShare pushErr = ", pushErr.Error())
	}
	base.DBEngine.Table("app_message").InsertOne(&appMessageOne)

	this.ReturnData = "success"
}

// @Title 后台查看调查问卷统计结果（h5使用）
// @Description 后台查看调查问卷统计结果（h5使用）
// @Param	sex					query			string	  		false		"性别"
// @Param	phoneNumber			query			string	  		false		"手机号"
// @Param	hasCourse			query			int		  		false		"是否领取课程，1是0否，-1或不传表示不筛选"
// @Param	pageNum				query 	  		int				true		"page num start from 1"
// @Param	pageTime			query 	  		int64			false		"page time should be empty when pagenum == 1"
// @Param	pageSize			query 	  		int				false		"page size default is 15"
// @Success 200 {object} models.QuestionnaireAuthorStatisticsForBackendContainer
// @router /getQuestionnaireAuthorStatisticsList [get]
func (this *QuestionnaireController) GetQuestionnaireAuthorStatisticsList() {
	sex := this.GetString("sex")
	phoneNumber := this.GetString("phoneNumber")
	hasCourse, _ := this.GetInt("hasCourse", -1)
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", util.UnixOfBeijingTime())
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}

	var freeCourseConfig models.SystemConfig
	base.DBEngine.Table("system_config").Where("program='questionnaire_free_course_id'").Get(&freeCourseConfig)

	var questionnaireAuthorStatisticsList []models.QuestionnaireAuthorStatisticsWithHasCourse
	totalSql := "select count(1) from questionnaire_author_statistics left join author on questionnaire_author_statistics.author_id=author.id where questionnaire_author_statistics.author_id != 0 and author.deleted_at is null and author.is_test=0 and questionnaire_author_statistics.author_id != 56 and questionnaire_author_statistics.deleted_at is null "
	dataSql := "select questionnaire_author_statistics.*, (select 1 from pay_order where pay_order.payer_id=questionnaire_author_statistics.author_id and pay_order.order_type=3 and pay_order.order_type_id='" + freeCourseConfig.ProgramValue + "' and pay_order.status=1 and pay_order.amount=0) as has_course from questionnaire_author_statistics left join author on questionnaire_author_statistics.author_id=author.id where questionnaire_author_statistics.author_id != 0 and author.deleted_at is null and author.is_test=0 and questionnaire_author_statistics.author_id != 56 and questionnaire_author_statistics.deleted_at is null "

	if sex != "" {
		totalSql += " and questionnaire_author_statistics.sex='" + sex + "' "
		dataSql += " and questionnaire_author_statistics.sex='" + sex + "' "
	}

	if phoneNumber != "" {
		totalSql += " and questionnaire_author_statistics.phone_number like '%" + phoneNumber + "%' "
		dataSql += " and questionnaire_author_statistics.phone_number like '%" + phoneNumber + "%' "
	}

	if hasCourse == 0 {
		totalSql += " and not exists (select 1 from pay_order where pay_order.payer_id=questionnaire_author_statistics.author_id and pay_order.status=1 and pay_order.amount=0 and pay_order.order_type=3 and pay_order.order_type_id='" + freeCourseConfig.ProgramValue + "' ) "
		dataSql += " and not exists (select 1 from pay_order where pay_order.payer_id=questionnaire_author_statistics.author_id and pay_order.status=1 and pay_order.amount=0 and pay_order.order_type=3 and pay_order.order_type_id='" + freeCourseConfig.ProgramValue + "' ) "
	} else if hasCourse == 1 {
		totalSql += " and exists (select 1 from pay_order where pay_order.payer_id=questionnaire_author_statistics.author_id and pay_order.status=1 and pay_order.amount=0 and pay_order.order_type=3 and pay_order.order_type_id='" + freeCourseConfig.ProgramValue + "' ) "
		dataSql += " and exists (select 1 from pay_order where pay_order.payer_id=questionnaire_author_statistics.author_id and pay_order.status=1 and pay_order.amount=0 and pay_order.order_type=3 and pay_order.order_type_id='" + freeCourseConfig.ProgramValue + "' ) "
	}

	dataSql += " order by questionnaire_author_statistics.created desc limit "+strconv.Itoa(pageSize*(pageNum-1))+", "+strconv.Itoa(pageSize)

	total, totalErr := base.DBEngine.SQL(totalSql).Count(new(models.QuestionnaireAuthorStatistics))
	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}

	if total > 0 {
		base.DBEngine.SQL(dataSql).Find(&questionnaireAuthorStatisticsList)
	}


	if questionnaireAuthorStatisticsList == nil {
		questionnaireAuthorStatisticsList = make([]models.QuestionnaireAuthorStatisticsWithHasCourse, 0)
	}

	this.ReturnData = models.QuestionnaireAuthorStatisticsForBackendContainer{models.BaseListContainer{total, pageNum, pageTime}, questionnaireAuthorStatisticsList}
}

// @Title 后台查看调查问卷详细回答（h5使用）
// @Description 后台查看调查问卷统计结果（h5使用）
// @Param	questionnaireAuthorStatisticsId		query	  		int64			true		"问卷统计id"
// @Success 200 {object} models.QuestionnaireAnswerListContainer
// @router /getQuestionnaireAnswerList [get]
func (this *QuestionnaireController) GetQuestionnaireAnswerList() {
	questionnaireAuthorStatisticsId := this.MustInt64("questionnaireAuthorStatisticsId")

	var questionnaireAuthorStatistics models.QuestionnaireAuthorStatistics
	hasStatistics, _ := base.DBEngine.Table("questionnaire_author_statistics").Where("id=?", questionnaireAuthorStatisticsId).Get(&questionnaireAuthorStatistics)
	if !hasStatistics {
		this.ReturnData = util.GenerateAlertMessage(models.QuestionnaireError700)
		return
	}

	var questionnaireAnswerList []models.QuestionnaireAnswer
	base.DBEngine.Table("questionnaire_answer").Where("questionnaire_author_statistics_id=?", questionnaireAuthorStatisticsId).Asc("created").Find(&questionnaireAnswerList)

	if questionnaireAnswerList == nil {
		questionnaireAnswerList = make([]models.QuestionnaireAnswer, 0)
	}

	this.ReturnData = models.QuestionnaireAnswerListContainer{questionnaireAnswerList}
}

//-------------------------------------------------方法-----------------------------------------------------------------------------------------------
//计算均分
func calcScoreAverage(intArr []int) float64 {
	sum := 0
	for _, val := range intArr {
		//累计求和
		sum += val
	}
	//平均值保留到2位小数
	return util.Rounding(float64(sum)/float64(len(intArr)))
}






