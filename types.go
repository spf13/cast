package cast

type (
	Float64er interface {
		Float64() (float64, error)
	}

	Int64er interface {
		Int64() (int64, error)
	}
)
