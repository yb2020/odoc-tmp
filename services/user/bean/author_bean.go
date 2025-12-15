package bean

// AuthorBean 作者信息
type AuthorBean struct {
	Id               string       `json:"id"`                // 用户Id
	NickName         string       `json:"nickName"`          // 昵称
	Description      string       `json:"description"`       // 描述
	Usn              string       `json:"usn"`               // 用户名
	UsnCanModify     bool         `json:"usnCanModify"`      // 用户名是否允许修改
	Mobile           string       `json:"mobile"`            // 手机
	Email            string       `json:"email"`             // 邮件
	IsWxPublicBind   bool         `json:"isWxPublicBind"`    // 公众号是否已绑定
	Self             bool         `json:"self"`              // 作者是否本人
	ShowName         string       `json:"showName"`          // 前端显示名称
	AvatarUrl        string       `json:"avatarUrl"`         // 头像URL
	Tags             string       `json:"tags"`              // 作者的标签
	AuthorId         string       `json:"authorId"`          // 作者id
	AuthorName       string       `json:"authorName"`        // 作者名称
	IsAuthentication bool         `json:"isAuthentication"`  // 是否已readpaper认证
	IsCert           bool         `json:"isCert"`            // 是否已认证
	IsPaperAuthor    bool         `json:"isPaperAuthor"`     // 是否是论文作者
	Profession       string       `json:"profession"`        // 职业
	ResearchField    string       `json:"researchField"`     // 研究领域
	SchoolCompany    string       `json:"schoolCompany"`     // 学校或公司
	BanInfo          *UserBanInfo `json:"banInfo,omitempty"` // 用户封禁信息
}

// UserBanInfo 用户封禁信息
type UserBanInfo struct {
	BanFlag    bool   `json:"banFlag"`    // 是否被封禁
	BanReason  string `json:"banReason"`  // 封禁原因
	BanRemark  string `json:"banRemark"`  // 封禁备注
	BanEndTime string `json:"banEndTime"` // 封禁结束时间
}
