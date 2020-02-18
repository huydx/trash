package trash

type Storage interface {
	Index(spans Spans) error
	GetTrace(id []byte) (Spans, error)
}

type BadgerStorage struct {

}

func (b BadgerStorage) Index(spans Spans) error {
	panic("implement me")
}

func (b BadgerStorage) GetTrace(id []byte) (Spans, error) {
	panic("implement me")
}

