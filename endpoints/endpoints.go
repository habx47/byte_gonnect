package endpoints

type Routes map[string]bool
type EndpointResps map[string][]byte

var HTTPVersion = "HTTP/1.1"
var SupportedMethods = []string{"GET", "POST", "PUT", "DELETE"}

var GETRsp_root = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 25\r\nConnection: close\r\n\r\nHello, from the root page")
var GETRsp_hello = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 26\r\nConnection: close\r\n\r\nHello, from the hello page")

var Rsp_NotFound = []byte("HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\nContent-Length: 14\r\nConnection: close\r\n\r\nPage not found")
var Rsp_MethodNotAllowed = []byte("HTTP/1.1 405 Method Not Allowed\r\nContent-Type: text/plain\r\nContent-Length: 18\r\nConnection: close\r\n\r\nMethod not allowed")
var Rsp_BadRequest = []byte("HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n")
