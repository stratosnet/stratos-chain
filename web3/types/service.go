package types

// Web3Service interface define service
type Web3Service interface{}

// Web3Context interface to define required elements for the API Context
type Web3Context interface {

	// service registry
	RegisterService(name string, srv Web3Service)
	ServiceList() map[string]Web3Service
}
