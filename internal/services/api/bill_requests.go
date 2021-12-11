package api

type BillRequests struct {
	url string
}

type BillRequestsI interface{}

func InitBillRequests(u string) BillRequestsI {
	return &BillRequests{
		url: u,
	}
}

// func (r *BillRequests) GetAll() (*fasthttp.Response, error)
