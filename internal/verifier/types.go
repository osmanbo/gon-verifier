package verifier

type (
	Request struct {
		TaskNo string
		User   UserInfo
		Params any
	}

	Response struct {
		TaskNo   string
		TeamName string
		Point    int32
		Reason   string
		Memo     string
	}

	Verifier interface {
		Do(req Request, res chan<- *Response)
		BuildParams(params [][]string) (any, error)
	}

	UserInfo struct {
		TeamName string
		Github   string
		Address  map[string]string
	}
)
