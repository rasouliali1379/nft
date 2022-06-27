package server
type IServer interface {
	ListenAndServe() error
	Shutdown() error
}
