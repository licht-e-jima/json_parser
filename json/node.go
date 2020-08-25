package calculation

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

	structNode struct {
		values [string]node
	}
)

func (*stringNode) dummy()
func (*numberNode) dummy()
func (*booleanNode) dummy()
func (*nullNode) dummy()
func (*arrayNode) dummy()
func (*structNode) dummy()
