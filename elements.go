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

type Margin struct {
	Amount Rect
	Elements
}

func NewMargin(amount Rect, els ...Element) *Margin {
	return &Margin{Amount: amount, Elements: Elements(els)}
}

func (margin *Margin) Draw(plot *Plot, canvas Canvas) {
	bounds := canvas.Bounds().Inset(margin.Amount)
	margin.Elements.Draw(plot, canvas.Context(bounds))
}

type VStack struct {
	Margin Rect
	Elements
}

func NewVStack(els ...Element) *VStack { return &VStack{Elements: Elements(els)} }

func (stack *VStack) Draw(plot *Plot, canvas Canvas) {
	if len(stack.Elements) == 0 {
		return
	}
	bounds := canvas.Bounds()
	for i, el := range stack.Elements {
		block := bounds.Row(i, len(stack.Elements))
		el.Draw(plot, canvas.Context(block.Inset(stack.Margin)))
	}
}

type HStack struct {
	Margin Rect
	Elements
}

func NewHStack(els ...Element) *HStack { return &HStack{Elements: Elements(els)} }

func (stack *HStack) Draw(plot *Plot, canvas Canvas) {
	if len(stack.Elements) == 0 {
		return
	}
	bounds := canvas.Bounds()
	for i, el := range stack.Elements {
		block := bounds.Column(i, len(stack.Elements))
		el.Draw(plot, canvas.Context(block.Inset(stack.Margin)))
	}
}
