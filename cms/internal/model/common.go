package model

// 通用响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(data interface{}) *Response {
	return &Response{Code: 0, Msg: "success", Data: data}
}

func Error(msg string) *Response {
	return &Response{Code: 1, Msg: msg}
}

func ErrorWithCode(code int, msg string) *Response {
	return &Response{Code: code, Msg: msg}
}

// 分页请求
type PageRequest struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

func (p *PageRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return (p.Page - 1) * p.PageSize
}

// 分页响应
type PageResponse struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     interface{} `json:"list"`
}

func NewPageResponse(total int64, page, pageSize int, list interface{}) *PageResponse {
	return &PageResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}
}
