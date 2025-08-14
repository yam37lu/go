package utils

import (
	"fmt"
	"time"
)

type FormattedTime time.Time

func (t FormattedTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}
