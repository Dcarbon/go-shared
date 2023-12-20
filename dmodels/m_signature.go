package dmodels

type Signature struct {
	Signer EthAddress `json:"signer"`
	Signed string     `json:"signed"` // Base64
	Data   string     `json:"data"`   // Base64
}

func (s *Signature) Verify() error {
	return s.Signer.VerifyPersonalSignBase64(s.Data, s.Signed)
}
