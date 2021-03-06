package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	types "github.com/binRick/go-gops-experiments/types"

	"github.com/google/gops/signal"
	"github.com/k0kubun/pp"
)

const (
	GOPS_HOST = `127.0.0.1`
	GOPS_PORT = 57223
)

var (
	GOPS_ADDR    = fmt.Sprintf(`%s:%d`, GOPS_HOST, GOPS_PORT)
	gops_signals = map[string]byte{
		"Stats":    signal.Stats,
		"Version":  signal.Version,
		"MemStats": signal.MemStats,
	}
	IGNORED_PROPERTIES = PropertiesList{
		"debug-gc",
		"enable-gc",
	}
)

type PropertiesList []string

func (list PropertiesList) Has(a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	gops_stats := types.GopsStats{
		Started: time.Now(),
		Success: false,
	}
	for label, s := range gops_signals {
		c, err := net.Dial("tcp", GOPS_ADDR)
		if err == nil {
			gops_stats.Success = true
			gops_stats.ConnectionsQty += 1
			buf := []byte{s}
			_, err = c.Write(buf)
			if err == nil {
				m, err := ioutil.ReadAll(c)
				if err == nil && len(m) > 0 {
					for _, l := range strings.Split(fmt.Sprintf(`%s`, m), "\n") {
						val_unit := `unknown`
						int_val := int64(0)
						if len(l) > 0 && len(strings.Split(l, `:`)) == 2 {
							prop := strings.TrimSpace(strings.Split(l, `:`)[0])
							val := strings.TrimSpace(strings.Split(l, `:`)[1])
							if !IGNORED_PROPERTIES.Has(prop) {
								if prop == `GOMAXPROCS` {
									b, int_err := strconv.ParseInt(val, 10, 0)
									if int_err == nil {
										val_unit = `procs`
										int_val = b
									}
								} else if prop == `num CPU` {
									b, int_err := strconv.ParseInt(val, 10, 0)
									if int_err == nil {
										val_unit = `cpus`
										int_val = b
									}
								} else if strings.HasSuffix(prop, " threads") {
									b, int_err := strconv.ParseInt(val, 10, 0)
									if int_err == nil {
										val_unit = `threads`
										int_val = b
									}
								} else if strings.HasSuffix(val, " bytes)") {
									l := strings.Split(val, `(`)
									b, int_err := strconv.ParseInt(strings.Split(l[len(l)-1], " bytes)")[0], 10, 0)
									if int_err == nil {
										val_unit = `bytes`
										int_val = b
									}
								}
								gops_stats.Stats = append(gops_stats.Stats, types.GopsStat{
									Property:  prop,
									Value:     val,
									ValueUnit: val_unit,
									IntValue:  int_val,
									Label:     label,
								})
							}
						}
					}
				}
			}
			c.Close()
			gops_stats.Duration = time.Since(gops_stats.Started)
		}
	}
	if false {
		pp.Println(gops_stats)
	}
	j, _ := json.Marshal(gops_stats)
	fmt.Println(string(j))
}
