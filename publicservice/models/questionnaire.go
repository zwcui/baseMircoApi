/*
@Time : 2019/12/16 下午12:00 
@Author : zwcui
@Software: GoLand
*/
package models

type Questionnaire struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	Name					string				`description:"更好人生学院学员信息统计表" json:"name"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type QuestionnaireQuestion struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	QuestionnaireId       	int64				`description:"用户id" json:"authorId"`
	Direction				string				`description:"测评方向" json:"direction"`
	Content					string				`description:"内容" json:"content"`
	Answer1					string				`description:"回答1" json:"answer1"`
	Score1					int					`description:"分数1" json:"score1"`
	Answer2					string				`description:"回答2" json:"answer2"`
	Score2					int					`description:"分数2" json:"score2"`
	Answer3					string				`description:"回答3" json:"answer3"`
	Score3					int					`description:"分数3" json:"score3"`
	Answer4					string				`description:"回答4" json:"answer4"`
	Score4					int					`description:"分数4" json:"score4"`
	Answer5					string				`description:"回答5" json:"answer5"`
	Score5					int					`description:"分数5" json:"score5"`
	DirectionWeight			float64				`description:"测评方向权重" json:"directionWeight"`
	SortNo					int					`description:"排序号" json:"sortNo"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type QuestionnaireAnswer struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	QuestionnaireQuestionId int64				`description:"问卷问题id" json:"questionnaireQuestionId"`
	QuestionnaireAuthorStatisticsId int64		`description:"问卷问题id" json:"questionnaireAuthorStatisticsId"`
	AuthorId       			int64				`description:"用户id" json:"authorId"`
	Direction				string				`description:"测评方向" json:"direction"`
	Content					string				`description:"内容" json:"content"`
	Answer					string				`description:"回答" json:"answer"`
	Score					int					`description:"分数" json:"score"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type QuestionnaireAuthorStatistics struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	QuestionnaireId 		int64				`description:"问卷id" json:"questionnaireId"`
	AuthorId       			int64				`description:"用户id" json:"authorId"`
	Name					string				`description:"姓名" json:"name"`
	Sex						string				`description:"性别" json:"sex"`
	Age						string				`description:"年龄" json:"age"`
	Occupation				string				`description:"职业" json:"occupation"`
	City					string				`description:"城市" json:"city"`
	Income					string				`description:"家庭年收入" json:"income"`
	Marriage				string				`description:"婚姻" json:"marriage"`
	IdCard					string				`description:"身份证号码" json:"idCard"`
	PhoneNumber				string				`description:"手机号" json:"phoneNumber"`
	Education				string				`description:"学历" json:"education"`
	Direction1				string				`description:"测评方向1" json:"direction1"`
	Average1				float64				`description:"均分1" json:"average1"`
	Direction2				string				`description:"测评方向2" json:"direction2"`
	Average2				float64				`description:"均分2" json:"average2"`
	Direction3				string				`description:"测评方向3" json:"direction3"`
	Average3				float64				`description:"均分3" json:"average3"`
	Direction4				string				`description:"测评方向4" json:"direction4"`
	Average4				float64				`description:"均分4" json:"average4"`
	Direction5				string				`description:"测评方向5" json:"direction5"`
	Average5				float64				`description:"均分5" json:"average5"`
	Direction6				string				`description:"测评方向6" json:"direction6"`
	Average6				float64				`description:"均分6" json:"average6"`
	Direction7				string				`description:"测评方向7" json:"direction7"`
	Average7				float64				`description:"均分7" json:"average7"`
	TotalScore				float64				`description:"分数" json:"totalScore"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//-------------------------------------结构体-----------------------------------------------------------------------
type QuestionnaireAnswerRequest struct {
	QuestionnaireQuestionId int64				`description:"问卷问题id" json:"questionnaireQuestionId"`
	AuthorId       			int64				`description:"用户id" json:"authorId"`
	Direction				string				`description:"测评方向" json:"direction"`
	Content					string				`description:"内容" json:"content"`
	Answer					string				`description:"回答" json:"answer"`
	Score					int					`description:"分数" json:"score"`
}


type QuestionnaireQuestionListContainer struct {
	QuestionnaireQuestionList	[]QuestionnaireQuestion	`description:"问题列表" json:"questionnaireQuestionList"`
}

type QuestionnaireAuthorStatisticsContainer struct {
	QuestionnaireAuthorStatistics	QuestionnaireAuthorStatistics	`description:"问卷结果" json:"questionnaireAuthorStatistics"`
	CourseName                      string				            `description:"赠送课程名称" json:"courseName"`
	CourseId                        int64				            `description:"课程id" json:"courseId"`
}

type QuestionnaireAuthorStatisticsWithHasCourse struct {
	QuestionnaireAuthorStatistics	`description:"问卷结果" json:"questionnaireAuthorStatistics" xorm:"extends"`
	HasCourse                      	int					            `description:"是否领取课程，1是0否" json:"hasCourse"`
}

type QuestionnaireAuthorStatisticsForBackendContainer struct {
	BaseListContainer
	QuestionnaireAuthorStatisticsList			[]QuestionnaireAuthorStatisticsWithHasCourse		`description:"问卷结果" json:"questionnaireAuthorStatisticsList"`
}

type QuestionnaireAnswerListContainer struct {
	QuestionnaireAnswerList			[]QuestionnaireAnswer			`description:"回答列表" json:"questionnaireAnswerList"`
}