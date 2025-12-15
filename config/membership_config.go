package config

type SubInfo struct {
	Type           uint32 `json:"type" yaml:"type"` // 订阅套餐类型, 对应constant.OrderType 0:Unknown, 1:Free, 2:PRO, 3:PRO_AddOnCredit
	Name           string `json:"name" yaml:"name"`
	Price          int64  `json:"price" yaml:"price"`
	OriginalPrice  int64  `json:"originalPrice" yaml:"originalPrice"`
	Currency       string `json:"currency" yaml:"currency"`
	Credit         int64  `json:"credit" yaml:"credit"`
	OriginalCredit int64  `json:"originalCredit" yaml:"originalCredit"`
	AddOnCredit    int64  `json:"addOnCredit" yaml:"addOnCredit"`
	Duration       int    `json:"duration" yaml:"duration"`
	StripePayMode  string `json:"stripePayMode" yaml:"stripePayMode"` // stripe支付模式 payment: 一次性付款, subscription: 订阅
	StripePriceId  string `json:"stripePriceId" yaml:"stripePriceId"` // stripe价格ID
}

type BasePermission struct {
	SubInfo                       SubInfo `json:"subInfo" yaml:"subInfo"`
	IsEnableAddOnCredit           bool    `json:"isEnableAddOnCredit" yaml:"isEnableAddOnCredit"`
	IsEnableSubAddOnCredit        bool    `json:"isEnableSubAddOnCredit" yaml:"isEnableSubAddOnCredit"`
	MaxAddOnCreditSubCountOfMonth int     `json:"maxAddOnCreditSubCountOfMonth" yaml:"maxAddOnCreditSubCountOfMonth"`
	SubAddOnCreditInfo            SubInfo `json:"subAddOnCreditInfo" yaml:"subAddOnCreditInfo"`
}

// Docs 文档权限配置
type DocsPermission struct {
	MaxStorageCapacity            int64 `json:"maxStorageCapacity" yaml:"maxStorageCapacity"`
	MaxStorageCapacityOriginal    int64 `json:"maxStorageCapacityOriginal" yaml:"maxStorageCapacityOriginal"`
	DocUploadMaxSize              int64 `json:"docUploadMaxSize" yaml:"docUploadMaxSize"`
	DocUploadMaxSizeOriginal      int64 `json:"docUploadMaxSizeOriginal" yaml:"docUploadMaxSizeOriginal"`
	DocUploadMaxPageCount         int32 `json:"docUploadMaxPageCount" yaml:"docUploadMaxPageCount"`
	DocUploadMaxPageCountOriginal int32 `json:"docUploadMaxPageCountOriginal" yaml:"docUploadMaxPageCountOriginal"`
}

// Note 笔记权限配置
type NotePermission struct {
	IsNoteSummary     bool `json:"isNoteSummary" yaml:"isNoteSummary"`
	IsNoteWord        bool `json:"isNoteWord" yaml:"isNoteWord"`
	IsNoteExtract     bool `json:"isNoteExtract" yaml:"isNoteExtract"`
	IsNoteManage      bool `json:"isNoteManage" yaml:"isNoteManage"`
	IsNotePdfDownload bool `json:"isNotePdfDownload" yaml:"isNotePdfDownload"`
}

// AI AI权限配置
type AiPermission struct {
	Copilot struct {
		IsEnable bool `json:"isEnable" yaml:"isEnable"`
		Models   []struct {
			Key        string `json:"key" yaml:"key"`
			Name       string `json:"name" yaml:"name"`
			IsEnable   bool   `json:"isEnable" yaml:"isEnable"`
			IsFree     bool   `json:"isFree" yaml:"isFree"`
			CreditCost int64  `json:"creditCost" yaml:"creditCost"`
		} `json:"models" yaml:"models"`
	} `json:"copilot" yaml:"copilot"`
}

// Translate 翻译权限配置
type TranslatePermission struct {
	IsOcr                               bool  `json:"isOcr" yaml:"isOcr"`
	OcrCreditCost                       int64 `json:"ocrCreditCost" yaml:"ocrCreditCost"`
	IsWordTranslate                     bool  `json:"isWordTranslate" yaml:"isWordTranslate"`
	WordTranslateCreditCost             int64 `json:"wordTranslateCreditCost" yaml:"wordTranslateCreditCost"`
	IsFullTextTranslate                 bool  `json:"isFullTextTranslate" yaml:"isFullTextTranslate"`
	FullTextTranslateCreditCost         int64 `json:"fullTextTranslateCreditCost" yaml:"fullTextTranslateCreditCost"`
	FullTextTranslateCreditCostOriginal int64 `json:"fullTextTranslateCreditCostOriginal" yaml:"fullTextTranslateCreditCostOriginal"`
	FullTextTranslateMaxPageCount       int64 `json:"fullTextTranslateMaxPageCount" yaml:"fullTextTranslateMaxPageCount"`
	IsAiTranslation                     bool  `json:"isAiTranslation" yaml:"isAiTranslation"`
	AiTranslationCreditCost             int64 `json:"aiTranslationCreditCost" yaml:"aiTranslationCreditCost"`
}

// MembershipTypeConfig 会员类型配置
type MembershipTypeConfig struct {
	Type        int32               `json:"type" yaml:"type"` // 会员类型，对应dto.MembershipType 0:Unknown, 1:Free, 2:PRO
	Name        string              `json:"name" yaml:"name"`
	Description string              `json:"description" yaml:"description"`
	IsFree      bool                `json:"isFree" yaml:"isFree"`
	Base        BasePermission      `json:"base" yaml:"base"`
	Docs        DocsPermission      `json:"docs" yaml:"docs"`
	Note        NotePermission      `json:"note" yaml:"note"`
	AI          AiPermission        `json:"ai" yaml:"ai"`
	Translate   TranslatePermission `json:"translate" yaml:"translate"`
}

// MembershipConfig 会员配置
type MembershipConfig struct {
	Free         MembershipTypeConfig `json:"free" yaml:"free"`
	Professional MembershipTypeConfig `json:"professional" yaml:"professional"`
}
