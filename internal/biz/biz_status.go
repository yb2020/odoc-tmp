package biz

// 系统状态码定义
// 状态码规则：
// - 1000-1999: 用户服务状态码
// - 2000-2999: 翻译服务状态码
// - 3000-3999: OAuth2服务状态码
// - 4000-4999: 会员服务状态码 （4000-4099: 用户会员/会员账号状态码，4100-4199: 会员订阅/订单状态码，4200-4699: 积分服务状态码）

const (
	// 用户服务状态码 (1000-1999)
	User_StatusNotFound        = 1001 // 用户不存在
	User_StatusAlreadyExists   = 1002 // 用户已存在
	User_StatusInvalidUserData = 1003 // 无效的用户数据
	User_StatusAuthFailed      = 1004 // 用户认证失败
	User_RSAKeyNotConfigured   = 1005 // RSA公钥未配置
	User_StatusCreated         = 1100 // 用户创建成功
	User_StatusUpdated         = 1101 // 用户更新成功
	User_StatusDeleted         = 1102 // 用户删除成功
	User_StatusCreateFailed    = 1103 // 创建用户失败
	User_StatusUpdateFailed    = 1104 // 更新用户失败
	User_StatusDeleteFailed    = 1105 // 删除用户失败
	User_StatusNotActive       = 1106 // 用户未激活

	// 翻译服务状态码 (2000-2999)
	Translate_StatusInvalidRequest  = 2006 // 无效的请求数据
	Translate_StatusNotFound        = 2007 // 翻译条目不存在
	Translate_StatusAlreadyExists   = 2008 // 翻译条目已存在
	Translate_StatusOperationFailed = 2009 // 操作失败
	Translate_StatusAllFailed       = 2010 // 所有翻译渠道都失败，内部翻译接口

	// OAuth2服务状态码 (3000-3999)
	OAuth2_StatusInvalidRequest       = 3001 // 无效的请求数据
	OAuth2_StatusInvalidCredentials   = 3002 // 无效的凭证
	OAuth2_StatusInvalidToken         = 3003 // 无效的令牌
	OAuth2_StatusTokenExpired         = 3004 // 令牌已过期
	OAuth2_StatusTokenRevoked         = 3005 // 令牌已撤销
	OAuth2_StatusTokenSigningFailed   = 3006 // 令牌签名失败
	OAuth2_StatusTokenStorageFailed   = 3007 // 令牌存储失败
	OAuth2_StatusTokenRetrievalFailed = 3008 // 令牌获取失败
	OAuth2_StatusUserNotFound         = 3009 // 用户不存在
	OAuth2_StatusOperationFailed      = 3010 // 操作失败

	// 会员服务状态码 (4000-4999)
	// 用户会员/会员账号状态码 (4000-4099)
	Membership_Status_UserAccountNotFound            = 4001 // 会员账户不存在
	Membership_Status_UserAccountAlreadyExists       = 4002 // 会员账户已存在
	Membership_Status_UserAccountNewAccountFailed    = 4003 // 会员账户创建失败
	Membership_Status_UserCreditAccountNotFound      = 4004 // 会员积分账户不存在
	Membership_Status_UserCreditAccountAlreadyExists = 4005 // 会员积分账户已存在
	Membership_Status_CreditBillTypeUnknown          = 4006 // 未知积分流水类型
	Membership_Status_CreditNotEnough                = 4007 // 积分不足
	Membership_Status_CreditAddOnNotEnough           = 4008 // 附加积分不足

	//会员订阅/订单状态码 (4100-4199)
	Membership_Status_SubscribeTypeNotFound                            = 4101 // 订阅类型不存在
	Membership_Status_CanNotSubscribeFree                              = 4102 // 不能订阅Free会员
	Membership_Status_CanNotSubscribePro                               = 4103 // 不能订阅Pro会员
	Membership_Status_CanNotSubscribeProAddOnCredit                    = 4104 // 不能订阅Pro会员附加积分
	Membership_Status_OverMaxAddOnCreditSubCountOfMonth                = 4105 // 超出Pro会员附加积分订阅次数
	Membership_Status_CreditPaymentRecordNotFound                      = 4106 // 支付记录不存在
	Membership_Status_CreditPaymentRecordStatusNotPending              = 4107 // 支付记录状态不是待支付
	Membership_Status_CreditPaymentRecordStatusNotAwaitingConfirmation = 4108 // 支付记录状态不是待确认

	//积分服务状态码 (4200-4699)
	Membership_Status_CreditServiceCheckPermissionAndPayConfigFunNoNull = 4201 // checkPermissionAndPayConfigFun 为空
	Membership_Status_CreditServicePermissionDenied                     = 4202 // 功能权限不足
	Membership_Status_CreditServiceTypeUnknown                          = 4203 // 积分服务类型未知错误

	Membership_Status_CreditService_Docs_Upload_OverFileSize   = 4301 // 积分服务-文档上传-文件大小超过限制
	Membership_Status_CreditService_Docs_Upload_OverPageSize   = 4302 // 积分服务-文档上传-页数超过限制
	Membership_Status_CreditService_Docs_Upload_OverMaxStorage = 4303 // 积分服务-文档上传-存储容量超过限制

	Membership_Status_CreditService_Ai_Copilot_NotEnabled      = 4401 // 积分服务-AI辅读-未开启
	Membership_Status_CreditService_Ai_Copilot_ModelNotFound   = 4402 // 积分服务-AI辅读-模型未找到
	Membership_Status_CreditService_Ai_Copilot_ModelNotEnabled = 4403 // 积分服务-AI辅读-模型未开启

	Membership_Status_CreditService_Translate_OcrNotEnabled         = 4501 // 积分服务-翻译-OCR未开启
	Membership_Status_CreditService_Translate_WordNotEnabled        = 4502 // 积分服务-翻译-划词翻译未开启
	Membership_Status_CreditService_Translate_FullTextNotEnabled    = 4503 // 积分服务-翻译-全文翻译未开启
	Membership_Status_CreditService_Translate_AiNotEnabled          = 4504 // 积分服务-翻译-AI翻译未开启
	Membership_Status_CreditService_Translate_FullText_OverPageSize = 4505 // 积分服务-翻译-全文翻译页数超过限制

	Membership_Status_CreditService_Note_SummaryNotEnabled     = 4601 // 积分服务-笔记-笔记总结未开启
	Membership_Status_CreditService_Note_WordNotEnabled        = 4602 // 积分服务-笔记-笔记单词未开启
	Membership_Status_CreditService_Note_ExtractNotEnabled     = 4603 // 积分服务-笔记-笔记摘录未开启
	Membership_Status_CreditService_Note_ManageNotEnabled      = 4604 // 积分服务-笔记-笔记管理未开启
	Membership_Status_CreditService_Note_PdfDownloadNotEnabled = 4605 // 积分服务-笔记-笔记pdf下载未开启

)
