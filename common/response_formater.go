package common

type Meta struct {
	Message string
	Status string
	Code int
}

type Response struct {
	Meta Meta
	Data interface{}
}

type ResponseList struct {
	Meta Meta
	Data interface {}
	Total int
}


func ResponseAPI(message string, status string, code int, data interface{} ) Response {
	meta := Meta {
		Message: message,
		Status: status,
		Code: code,
	}

	response := Response {
		Meta: meta, 
		Data: data,
	}

	return response
}

func ResponseAPIList(message string, status string, code int, data interface{}, total int) ResponseList {
	meta := Meta {
		Message: message,
		Status: status,
		Code: code,
	}

	responseList := ResponseList {
		Meta: meta, 
		Data: data,
		Total: total,
	}

	return responseList
}
