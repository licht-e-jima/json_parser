package json

type (
	node interface {
		dummy()
	}

	stringNode struct {
		char string
	}

	numberNode struct {
		char string
	}

	booleanNode struct {
		boolean bool
	}

	nullNode struct{}

	arrayNode struct {
		elements []node
	}

	structureNode struct {
		values map[string]node
	}
)

func (*stringNode) dummy()    {}
func (*numberNode) dummy()    {}
func (*booleanNode) dummy()   {}
func (*nullNode) dummy()      {}
func (*arrayNode) dummy()     {}
func (*structureNode) dummy() {}
