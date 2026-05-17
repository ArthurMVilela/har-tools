package model

import (
	"net"
	"time"
)

type HAR struct {
	Log Log `json:"log"` // Root of the exported data
}

type Log struct {
	Version string   `json:"version"`           // Version number of the format
	Creator Creator  `json:"name"`              // Information about the log creator application
	Browser *Browser `json:"browser"`           // Information of the user agent
	Pages   []Page   `json:"pages"`             // Exported (tracked) pages, if the application supports grouping by page
	Entries []Entry  `json:"entries"`           // Exported (tracked) HTTP request
	Comment string   `json:"comment,omitempty"` // Comment provided by the user or application
}

type Creator struct {
	Name    string `json:"name"`              // The name of the application that created the log
	Version string `json:"version"`           // The version number of the application that created the log
	Comment string `json:"comment,omitempty"` // A comment provided by the user or the application
}

type Browser struct {
	Name    string `json:"name"`              // The name of the browser that created the log
	Version string `json:"version"`           // The version number of the browser that created the log
	Comment string `json:"comment,omitempty"` // A comment provided by the user or the browser
}

type Page struct {
	ID              string      `json:"id"`                // Unique identifier of the page
	StartedDateTime time.Time   `json:"startedDateTime"`   // Date and time stamp for the beginning of the page load
	Title           string      `json:"title"`             // Page title
	PageTimings     PageTimings `json:"pageTimings"`       // Detailed timing info about page load
	Comment         string      `json:"comment,omitempty"` // A comment provided by the user or the application
}

type PageTimings struct {
	OnContentLoad int    `json:"onContentLoad,omitempty"` // Content of the page loaded. Number of milliseconds since page load started
	OnLoad        int    `json:"onLoad,omitempty"`        // Page is loaded (onLoad event fired). Number of milliseconds since page load started (page.startedDateTime)
	Comment       string `json:"comment,omitempty"`       // A comment provided by the user or the application
}

type Entry struct {
	Pageref         string     `json:"pageref,omitempty"`        // Reference to parent page
	StartedDateTime time.Time  `json:"startedDateTime"`          // Date and time stamp of the request start
	Time            int        `json:"time"`                     // Total elapsed time of the request in milliseconds
	Request         Request    `json:"request"`                  // Detailed info about the request
	Response        Response   `json:"response"`                 // Detailed info about the response
	Cache           Cache      `json:"cache"`                    // Info about cache usage
	Timings         Timings    `json:"timings"`                  // Detailed timing info about request/response round trip
	ServerIPAddress net.IPAddr `json:"serverIPAddress,omitzero"` // IP address of the server that was connected (result of DNS resolution)
	Connection      string     `json:"connection,omitempty"`     // Unique ID of the parent TCP/IP connection, can be the client port number. Note that a port number doesn't have to be unique identifier in cases where the port is shared for more connections. If the port isn't available for the application, any other unique connection ID can be used instead (e.g. connection index). Leave out this field if the application doesn't support this info.
	Comment         string     `json:"comment,omitempty"`        // A comment provided by the user or the application
}

type Request struct {
	Method      string        `json:"method"`             // Request method
	URL         string        `json:"url"`                // Absolute URL of the request
	HTTPVersion string        `json:"httpVersion"`        // Request HTTP Version
	Cookies     []Cookies     `json:"cookies"`            // List of cookie objects
	Headers     []Header      `json:"header"`             // List of header objects
	QueryString []QueryString `json:"queryString"`        // List of query parameter objects
	PostData    *PostData     `json:"postData,omitempty"` // Posted data info
	HeaderSize  int           `json:"headersSize"`        // Total number of bytes from the start of the HTTP request message until (and including) the double CRLF before the body
	BodySize    int           `json:"bodySize"`           // Size of the request body (POST data payload) in bytes
	Comment     string        `json:"comment,omitempty"`  // A comment provided by the user or the application
}

type Response struct {
	Status      int       `json:"status"`            // Response status
	StatusText  string    `json:"statusText"`        // Response status description
	HTTPVersion string    `json:"httpVersion"`       // Response HTTP Version
	Cookies     []Cookies `json:"cookies"`           // List of cookie objects
	Headers     []Header  `json:"header"`            // List of header objects
	Content     Content   `json:"content"`           // Details about the response body
	RedirectURL string    `json:"redirectURL "`      // Redirection target URL from the Location response header
	HeadersSize int       `json:"headersSize "`      // Total number of bytes from the start of the HTTP response message until (and including) the double CRLF before the body. Set to -1 if the info is not available
	BodySize    int       `json:"bodySize"`          // Size of the received response body in bytes. Set to zero in case of responses coming from the cache (304). Set to -1 if the info is not available.
	Comment     string    `json:"comment,omitempty"` // A comment provided by the user or the application
}

type Cookies struct {
	Name     string    `json:"name"`               // The name of the cookie
	Value    string    `json:"value"`              // The cookie value
	Path     string    `json:"path,omitempty"`     // The path pertaining to the cookie
	Domain   string    `json:"domain,omitempty"`   // The host of the cookie
	Expires  time.Time `json:"expires,omitzero"`   // Cookie expiration time
	HTTPOnly bool      `json:"httpOnly,omitempty"` // Set to true if the cookie is HTTP only, false otherwise
	Secure   bool      `json:"secure,omitempty"`   // True if the cookie was transmitted over ssl, false otherwise
	Comment  string    `json:"comment,omitempty"`  // A comment provided by the user or the application
}

type Header struct {
	Name    string `json:"name"`              // Header name
	Value   string `json:"value"`             // Header Value
	Comment string `json:"comment,omitempty"` // A comment provided by the user or the application
}

type QueryString struct {
	Name    string `json:"name"`              // Query name
	Value   string `json:"value"`             // Query Value
	Comment string `json:"comment,omitempty"` // A comment provided by the user or the application
}

type PostData struct {
	MimeType string      `json:"mimeType"`          // Mime type of posted data
	Params   []URLParams `json:"params"`            // List of posted parameters (in case of URL encoded parameters)
	Text     string      `json:"text"`              // Plain text posted data
	Comment  string      `json:"comment,omitempty"` // A comment provided by the user or the application
}

type URLParams struct {
	Name        string `json:"name"`                  // Name of a posted parameter
	Value       string `json:"value,omitempty"`       // Value of a posted parameter or content of a posted file
	FileName    string `json:"fileName,omitempty"`    // Name of a posted file
	ContentType string `json:"contentType,omitempty"` // Content type of a posted file
	Comment     string `json:"comment,omitempty"`     // A comment provided by the user or the application
}

type Content struct {
	Size        int    `json:"size"`                  // Length of the returned content in bytes. Should be equal to response.bodySize if there is no compression and bigger when the content has been compressed
	Compression int    `json:"compression,omitempty"` // Number of bytes saved. Leave out this field if the information is not available
	MimeType    string `json:"mimeType"`              // MIME type of the response text (value of the Content-Type response header). The charset attribute of the MIME type is included (if available)
	Text        string `json:"text,omitempty"`        // Response body sent from the server or loaded from the browser cache. This field is populated with textual content only. The text field is either HTTP decoded text or a encoded (e.g. "base64") representation of the response body. Leave out this field if the information is not available.
	Encoding    string `json:"encoding,omitempty"`    // Encoding used for response text field e.g "base64". Leave out this field if the text field is HTTP decoded (decompressed & unchunked), than trans-coded from its original character set into UTF-8.
	Comment     string `json:"comment,omitempty"`     // A comment provided by the user or the application
}

type Cache struct {
	BeforeRequest *CacheData `json:"beforeRequest,omitempty"` // State of a cache entry before the request. Leave out this field if the information is not available
	AfterRequest  *CacheData `json:"afterRequest,omitempty"`  // State of a cache entry after the request. Leave out this field if the information is not available
	Comment       string     `json:"comment,omitempty"`       // A comment provided by the user or the application
}

type CacheData struct {
	Expires    time.Time `json:"expires,omitzero"`  // Expiration time of the cache entry
	LastAccess time.Time `json:"lastAccess"`        // The last time the cache entry was opened
	ETag       string    `json:"eTag"`              // Etag
	HitCount   int       `json:"hitCount"`          // The number of times the cache entry has been opened
	Comment    string    `json:"comment,omitempty"` // A comment provided by the user or the application
}

// Timings describes various phases within request-response round trip. All times are specified in milliseconds.
type Timings struct {
	Blocked int `json:"blocked,omitempty"` // Time spent in a queue waiting for a network connection. Use -1 if the timing does not apply to the current request
	DNS     int `json:"dns,omitempty"`     // DNS resolution time. The time required to resolve a host name. Use -1 if the timing does not apply to the current request
	Connect int `json:"connect,omitempty"` // Time required to create TCP connection. Use -1 if the timing does not apply to the current request
	Send    int `json:"send"`              // Time required to send HTTP request to the server
	Wait    int `json:"wait"`              // Waiting for a response from the server
	Receive int `json:"receive"`           // Time required to read entire response from the server (or cache)
	// Time required for SSL/TLS negotiation. If this field is defined then the time is also included in the connect field
	// (to ensure backward compatibility with HAR 1.1). Use -1 if the timing does not apply to the current request.
	SSL     int    `json:"ssl,omitempty"`
	Comment string `json:"comment,omitempty"` // A comment provided by the user or the application
}
