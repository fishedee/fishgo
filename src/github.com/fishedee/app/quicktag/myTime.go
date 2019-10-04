package quicktag

import (
	"time"
)

type myTime time.Time

func (this myTime) MarshalJSON() ([]byte, error) {
	t := time.Time(this)
	return []byte(t.Format("\"2006-01-02 15:04:05\"")), nil
}

func (this *myTime) UnmarshalJSON(data []byte) error {
	t, err := time.ParseInLocation("\"2006-01-02 15:04:05\"", string(data), time.Local)
	if err != nil {
		return err
	}
	*this = myTime(t)
	return nil
}
