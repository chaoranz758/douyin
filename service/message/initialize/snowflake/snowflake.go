package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var Node *sf.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	Node, err = sf.NewNode(machineID)
	return
}
