package cast

type (
	float64er interface {
		Float64() (float64, error)
	}

	int64er interface {
		Int64() (int64, error)
	}
)
