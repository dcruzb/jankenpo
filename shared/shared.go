package shared

const SAMPLE_SIZE = 5
const SERVER_PORT = 46000

type Request struct {
	Player1 string
	Player2 string
}

type Reply struct {
	Result int
}
