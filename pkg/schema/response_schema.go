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

func (r *Response) addClusters(clusters []Cluster) {
	// for i, s := range clusters.GkeCluster {
	// 	cr:= ClusterResponse{
	// 		Name: s.Name,
	// 		Type: "GKE",
	// 	}

	// 	r.ClusterResponses = append(s.ClusterExpect, parroty.GlobalExpect.ClusterExpect...)
	// }

	// for i, s := range clusters.EksCluster {
	// 	parroty.Clusters.Cluster[i].ClusterExpect = append(s.ClusterExpect, parroty.GlobalExpect.ClusterExpect...)
	// }

	// for i, s := range clusters.Cluster {
	// 	parroty.Clusters.Cluster[i].ClusterExpect = append(s.ClusterExpect, parroty.GlobalExpect.ClusterExpect...)
	// }
}
