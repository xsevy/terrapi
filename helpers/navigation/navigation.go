package navigation

type NavigableItem interface {
	GetName() string
	GetChildren() []NavigableItem
	IsDisabled() bool
	GetID() string
}

type NavigationItem struct {
	Items    []NavigableItem
	Selected Selected
}

type NavigationStack struct {
	Stack []NavigationItem
}

func NewNavigationStack(initialItems []NavigableItem) NavigationStack {
	return NavigationStack{
		Stack: []NavigationItem{{Items: initialItems}},
	}
}

func (ns *NavigationStack) CurrentItem() *NavigationItem {
	if len(ns.Stack) == 0 {
		return nil
	}
	return &ns.Stack[len(ns.Stack)-1]
}

func (ns *NavigationStack) Push(item NavigableItem) {
	ns.Stack = append(ns.Stack, NavigationItem{Items: item.GetChildren()})
}

func (ns *NavigationStack) Pop() {
	if len(ns.Stack) > 1 {
		ns.Stack = ns.Stack[:len(ns.Stack)-1]
	}
}
