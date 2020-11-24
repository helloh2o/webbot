package bean

type Topic struct {
	Id         int64  `json:"topicId"`
	Type       int64  `json:"type"`
	Title      string `json:"title"`
	Views      int64  `json:"viewCount"`
	CreateTime string `json:"createTime"`
}

type Page struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type View struct {
	Topics []Topic
	Page   Page
}
