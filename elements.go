package plot

type Elements []Element

func (els *Elements) Add(el Element)           { *els = append(*els, el) }
func (els *Elements) AddGroup(adds ...Element) { els.Add(Elements(adds)) }

func (els Elements) Stats() Stats { return maximalStats(els) }
func (els Elements) Draw(plot *Plot, canvas Canvas) {
	for _, el := range els {
		el.Draw(plot, canvas)
	}
}

type VStack struct{ Elements }

func NewVStack(els ...Element) *VStack { return &VStack{Elements(els)} }

func (stack *VStack) Draw(plot *Plot, canvas Canvas) {
	if len(stack.Elements) == 0 {
		return
	}
	bounds := canvas.Bounds()
	for i, el := range stack.Elements {
		el.Draw(plot, canvas.Context(bounds.Row(i, len(stack.Elements))))
	}
}

type HStack struct{ Elements }

func NewHStack(els ...Element) *HStack { return &HStack{Elements(els)} }

func (stack *HStack) Draw(plot *Plot, canvas Canvas) {
	if len(stack.Elements) == 0 {
		return
	}
	bounds := canvas.Bounds()
	for i, el := range stack.Elements {
		el.Draw(plot, canvas.Context(bounds.Column(i, len(stack.Elements))))
	}
}
