package message

type Platform struct {
	Flash        string `json:"flash"`
	Download     string `json:"download"`
	HTML5Mobile  string `json:"html5mobile"`
	HTML5Desktop string `json:"html5desktop"`
	Nativemobile string `json:"nativemobile"`
}
