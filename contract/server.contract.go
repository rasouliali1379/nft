package contract

type IServer interface {
	ListenAndServe() error
	Shutdown() error
}
