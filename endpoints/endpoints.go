package endpoints

type Routes map[string]bool
type EndpointResps map[string][]byte

var ValidTokenChars = [256]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true,
	'f': true, 'g': true, 'h': true, 'i': true, 'j': true,
	'k': true, 'l': true, 'm': true, 'n': true, 'o': true,
	'p': true, 'q': true, 'r': true, 's': true, 't': true,
	'u': true, 'v': true, 'w': true, 'x': true, 'y': true,
	'z': true,
	'A': true, 'B': true, 'C': true, 'D': true, 'E': true,
	'F': true, 'G': true, 'H': true, 'I': true, 'J': true,
	'K': true, 'L': true, 'M': true, 'N': true, 'O': true,
	'P': true, 'Q': true, 'R': true, 'S': true, 'T': true,
	'U': true, 'V': true, 'W': true, 'X': true, 'Y': true,
	'Z': true,
	'!': true, '#': true, '$': true, '%': true, '&': true,
	'\'': true, '*': true, '+': true, '-': true, '.': true,
	'^': true, '_': true, '`': true, '|': true, '~': true,
}

var HTTPVersion = "HTTP/1.1"
var SupportedMethods = []string{"GET", "POST", "PUT", "DELETE"}

var GETRsp_root = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 25\r\nConnection: close\r\n\r\nHello, from the root page")
var GETRsp_hello = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 26\r\nConnection: close\r\n\r\nHello, from the hello page")

var Rsp_NotFound = []byte("HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\nContent-Length: 14\r\nConnection: close\r\n\r\nPage not found")
var Rsp_MethodNotAllowed = []byte("HTTP/1.1 405 Method Not Allowed\r\nContent-Type: text/plain\r\nContent-Length: 18\r\nConnection: close\r\n\r\nMethod not allowed")
var Rsp_BadRequest = []byte("HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n")
