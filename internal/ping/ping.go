package ping

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-ping/ping"
)

type Monitor struct {
	Destinations []string
	Database     *sql.DB
	Pingers      []*ping.Pinger
}

// NewMonitor creates a new Monitor instance
func NewMonitor(destinations []string, database *sql.DB) (*Monitor, error) {
	monitor := &Monitor{
		Destinations: destinations,
		Database:     database,
	}

	for _, destination := range destinations {
		pinger, err := ping.NewPinger(destination)
		if err != nil {
			return nil, err
		}

		pinger.Interval = time.Second * 5

		pinger.OnRecv = func(pkt *ping.Packet) {
			go monitor.updateDatabase(pkt)
			log.Println("Seq:", pkt.Seq, pkt.Nbytes, "bytes received from", pkt.Addr, "in", pkt.Rtt)
		}

		pinger.OnSend = func(pkt *ping.Packet) {
			log.Println("Seq:", pkt.Seq, pkt.Nbytes, "bytes sent to", pkt.Addr)
			go monitor.insertIntoDatabase(pkt)
		}

		monitor.Pingers = append(monitor.Pingers, pinger)
	}

	return monitor, nil
}

func (m *Monitor) Start() {
	for _, pinger := range m.Pingers {
		go m.runPinger(pinger)
	}
}

func (m *Monitor) Statistics() []*ping.Statistics {
	var stats []*ping.Statistics

	for _, pinger := range m.Pingers {
		stats = append(stats, pinger.Statistics())
	}

	return stats
}

func (m *Monitor) runPinger(pinger *ping.Pinger) {
	err := pinger.Run()

	if err != nil {
		log.Println("Ping error", err)
	}
}

func (m *Monitor) updateDatabase(pkt *ping.Packet) {
	_, err := m.Database.Exec("UPDATE monitoring SET status = ?, response_time = ? WHERE id = (SELECT MAX(id) FROM monitoring WHERE seq_id = ? AND source = ?);", 1, pkt.Rtt.Milliseconds(), pkt.Seq, pkt.Addr)

	if err != nil {
		log.Println("Error updating data in the database:", err)
	}
}

func (m *Monitor) insertIntoDatabase(pkt *ping.Packet) {
	_, err := m.Database.Exec("INSERT INTO monitoring (seq_id, source, status, response_time) VALUES (?, ?, ?, ?)", pkt.Seq, pkt.Addr, 0, 0)

	if err != nil {
		log.Println("Error inserting data into the database:", err)
	}
}
