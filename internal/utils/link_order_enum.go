package utils

type LinkOrderBy int

const (
	OrderByCreatedDate LinkOrderBy = iota
	OrderByUpdatedDate
	OrderByExpiredDate
	OrderByCounts
)

func (l LinkOrderBy) GetString() string {
	switch l {
	case OrderByCounts:
		return "counts"
	case OrderByCreatedDate:
		return "created_at"
	case OrderByUpdatedDate:
		return "updated_at"
	case OrderByExpiredDate:
		return "expired_at"
	}

	return "created_at"
}

func ParseLinkOrderBy(s string) (LinkOrderBy, error) {
	switch s {
	case "created_at":
		return OrderByCreatedDate, nil
	case "updated_at":
		return OrderByUpdatedDate, nil
	case "expired_at":
		return OrderByExpiredDate, nil
	case "counts":
		return OrderByCounts, nil
	default:
		return OrderByCreatedDate, &InvalidOrderByError{Value: s}
	}
}

func ParseLinkOrderByOrDefault(s string) LinkOrderBy {
	orderBy, _ := ParseLinkOrderBy(s)
	return orderBy
}

type InvalidOrderByError struct {
	Value string
}

func (e *InvalidOrderByError) Error() string {
	return "invalid order by value: " + e.Value + ". Valid values are: created_at, updated_at, expired_at, counts"
}
