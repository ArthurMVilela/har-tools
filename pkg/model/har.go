package model

import "time"

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
	Pageref         string    `json:"pageref,omitempty"` // Reference to parent page
	StartedDateTime time.Time `json:"startedDateTime"`   // Date and time stamp of the request start
	Time            int       `json:"time"`              // Total elapsed time of the request in milliseconds
	Request         Request   `json:"request"`           // Detailed info about the request
	Response        Response  `json:"response"`          // Detailed info about the response
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
