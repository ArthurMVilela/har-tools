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
	Request         any       `json:"request"`           // Detailed info about the request
	Response        any       `json:"response"`          // Detailed info about the response
}
