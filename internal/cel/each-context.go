package cel

type EachContext struct {
	Index int
	Item  any
}

func NewEachContext() EachContext {
	return EachContext{
		Index: 0,
		Item:  nil,
	}
}

func (e *EachContext) GetItem() any {
	return e.Item
}

func (e *EachContext) GetItemIndex() int {
	return e.Index
}

func (e *EachContext) SetItem(item any) {
	e.Item = item
}

func (e *EachContext) SetItemIndex(index int) {
	e.Index = index
}
