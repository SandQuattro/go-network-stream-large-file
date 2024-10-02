package proto

const (
	TCP = iota
	UDP
)

const MaxPacketSize = 1000

func String(proto int) string {
	switch proto {
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	default:
		return "error"
	}
}
