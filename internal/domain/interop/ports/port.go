package ports

type Facade interface {
	isFacade()
}

type BaseFacade struct{}

func (BaseFacade) isFacade() {}
