package ping

import (
	"database/sql"
	"log"

	"github.com/go-ping/ping"
)

type Monitor struct {
	Destinations []string
	Database     *sql.DB
	pingers      []*ping.Pinger
}

// NewMonitor creates a new Monitor instance
func NewMonitor(destinations []string, database *sql.DB) (monitor *Monitor, err error) {
	monitor = &Monitor{
		Destinations: destinations,
		Database:     database,
	}

	for _, destination := range destinations {
		pinger, err := ping.NewPinger(destination)
		if err != nil {
			return nil, err
		}

		pinger.OnRecv = func(pkt *ping.Packet) {
			_, err = database.Exec("UPDATE monitoring SET status = ?, response_time = ? WHERE seq_id = ? AND source = ?", 1, pkt.Rtt.Milliseconds(), pkt.Seq, pkt.Addr)

			if err != nil {
				log.Println("Error updating data in the database:", err)
			}

			log.Println("Seq: ", pkt.Seq, " ", pkt.Nbytes, "bytes received from", pkt.Addr, "in", pkt.Rtt)
		}

		pinger.OnSend = func(pkt *ping.Packet) {
			log.Println("Seq: ", pkt.Seq, " ", pkt.Nbytes, "bytes sent to", pkt.Addr)

			_, err = database.Exec("INSERT INTO monitoring (seq_id, source, status, response_time) VALUES (?, ?, ?, ?)", pkt.Seq, pkt.Addr, 0, 0)

			if err != nil {
				log.Println("Error inserting data into the database:", err)
			}
		}

		monitor.pingers = append(monitor.pingers, pinger)
	}

	return monitor, nil
}

func (m *Monitor) Start() {
	for _, pinger := range m.pingers {
		go runPinger(pinger)
	}
}

func (m *Monitor) Statistics() []*ping.Statistics {
	var statss []*ping.Statistics

	for _, pinger := range m.pingers {
		stats := pinger.Statistics()

		statss = append(statss, stats)
	}

	return statss
}

func runPinger(pinger *ping.Pinger) {
	err := pinger.Run()

	if err != nil {
		log.Println("ping error", err)
	}
}
