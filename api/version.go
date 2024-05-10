package api

type VersionInfoResponse struct {
	Name      string `json:"name"`
	BuildDate string `json:"buildDate"`
	GitHash   string `json:"gitHash"`
	Version   string `json:"version"`
}
