package schema

type Response struct {
	ClusterResponses []ClusterResponse
}

type ClusterResponse struct {
	Name   string
	Type   string
	Checks []Check
}

type Check struct {
	ClusterExpect
	Passed bool
}

func (r *Response) AddClusterResponse(cr ClusterResponse) {
	r.ClusterResponses = append(r.ClusterResponses, cr)
}

func (cr *ClusterResponse) AddCheck(ce ClusterExpect, passed bool) {
	c := Check{
		ClusterExpect: ce,
		Passed:        passed,
	}
	cr.Checks = append(cr.Checks, c)
}
