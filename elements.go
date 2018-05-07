package plot

type Elements []Element

func (els *Elements) Add(el Element)           { *els = append(*els, el) }
func (els *Elements) AddGroup(adds ...Element) { els.Add(Elements(adds)) }

func (els Elements) Stats() Stats { return maximalStats(els) }
func (els Elements) Draw(plot *Plot, canvas Canvas) {
	for _, el := range els {
		if el == nil {
			continue
		}
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

type HFlex struct {
	Margin Rect

	fixedSize []float64
	elements  Elements
}

func NewHFlex() *HFlex { return &HFlex{} }

func (stack *HFlex) Stats() Stats { return stack.elements.Stats() }

func (stack *HFlex) Add(fixedSize float64, el Element) {
	stack.elements.Add(el)
	stack.fixedSize = append(stack.fixedSize, fixedSize)
}
func (stack *HFlex) AddGroup(fixedSize float64, adds ...Element) {
	stack.Add(0, Elements(adds))
}

func (stack *HFlex) Draw(plot *Plot, canvas Canvas) {
	if len(stack.elements) == 0 {
		return
	}

	fixedSize := 0.0
	flexCount := 0.0
	for i, size := range stack.fixedSize {
		fixedSize += size
		if stack.elements[i] == nil {
			continue
		}
		if size == 0 {
			flexCount++
		}
	}

	bounds := canvas.Bounds()
	size := bounds.Size()

	flexWidth := (bounds.Size().X - fixedSize) / flexCount
	min := bounds.Min
	for i, el := range stack.elements {
		elsize := stack.fixedSize[i]
		if el == nil {
			min.X += elsize
			continue
		}
		if elsize == 0 {
			elsize = flexWidth
		}

		block := Rect{
			min,
			min.Add(Point{elsize, size.Y}),
		}
		min.X = block.Max.X

		el.Draw(plot, canvas.Context(block.Inset(stack.Margin)))
	}
}
