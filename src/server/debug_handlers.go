package main

import (
	"html/template"
	"net/http"
)

// WsDebugHandler is a webpage with a client for the normal websocket handler.
func WsDebugHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, &homeTemplateCtx{
		Host:         "ws://" + r.Host + "/ws",
		Title:        "Websocket debug client",
		DefaultValue: "{ \"message_type\": \"select\", \"select\": \"arne\" } ",
	})
}

// WsDebugControlPanelHandler is a webpage with a client for the control panel.
func WsDebugControlPanelHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, &homeTemplateCtx{
		Host:         "ws://" + r.Host + "/ws/control_panel",
		Title:        "Control panel websocket debug client",
		DefaultValue: "{ \"program\": 1, \"color\": [128, 255, 0], \"brightness\": 255 }",
	})
}

// WsDebugOpcHandler is a webpage with a client for the control panel.
func WsDebugOpcHandler(w http.ResponseWriter, r *http.Request) {
	opcTemplate.Execute(w, &homeTemplateCtx{
		Host:         "ws://" + r.Host + "/ws/opc",
		Title:        "OPC websocket debug client",
		DefaultValue: "",
	})
}

type homeTemplateCtx struct {
	Host         string
	Title        string
	DefaultValue string
}

var homeTemplate = template.Must(template.New("").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	<meta charset="utf-8">
	<script>  
	window.addEventListener("load", function(evt) {
		var output = document.getElementById("output");
		var input = document.getElementById("input");
		var ws;
		var print = function(message) {
			var d = document.createElement("div");
			d.innerHTML = message;
			output.appendChild(d);
		};
		document.getElementById("open").onclick = function(evt) {
			if (ws) {
				return false;
			}
			ws = new WebSocket("{{ .Host }}");
			ws.onopen = function(evt) {
				print("OPEN");
			}
			ws.onclose = function(evt) {
				print("CLOSE");
				ws = null;
			}
			ws.onmessage = function(evt) {
				print("RESPONSE: " + evt.data);
			}
			ws.onerror = function(evt) {
				print("ERROR: " + evt.data);
			}
			return false;
		};
		document.getElementById("send").onclick = function(evt) {
			if (!ws) {
				return false;
			}
			print("SEND: " + input.value);
			ws.send(input.value);
			return false;
		};
		document.getElementById("close").onclick = function(evt) {
			if (!ws) {
				return false;
			}
			ws.close();
			return false;
		};
	});
	</script>
	</head>
	<body>
	<table>
	<tr><td valign="top" width="50%">
	<h1>{{ .Title }}</h1>
	<p>Click "Open" to create a connection to the server, 
	"Send" to send a message to the server and "Close" to close the connection. 
	You can change the message and send multiple times.
	<p>
	<form>
	<button id="open">Open</button>
	<button id="close">Close</button>
	<p><input id="input" type="text" width="400px" value="{{ .DefaultValue }}">
	<button id="send">Send</button>
	</form>
	</td><td valign="top" width="50%">
	<div id="output"></div>
	</td></tr></table>
	</body>
	</html>
	`))

var opcTemplate = template.Must(template.New("").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
		<meta charset="utf-8">
		<script>  
		window.addEventListener("load", function(evt) {
			var output = document.getElementById("output");
			var input = document.getElementById("input");
			var ws;
			var print = function(message) {
				var d = document.createElement("div");
				d.innerHTML = message;
				output.appendChild(d);
			};
			document.getElementById("open").onclick = function(evt) {
				if (ws) {
					return false;
				}
				ws = new WebSocket("{{ .Host }}");
				ws.binaryType = 'arraybuffer';

				ws.onopen = function(evt) {
					print("OPEN");
				}
				ws.onclose = function(evt) {
					print("CLOSE");
					ws = null;
				}
				ws.onmessage = function(evt) {
					print("RESPONSE: " + evt.data);
					//var data = evt.data;
				}
				ws.onerror = function(evt) {
					print("ERROR: " + evt.data);
				}
				return false;
			};
			document.getElementById("send").onclick = function(evt) {
				if (!ws) {
					return false;
				}
				print("SEND: " + input.value);
				ws.send(input.value);
				return false;
			};
			document.getElementById("close").onclick = function(evt) {
				if (!ws) {
					return false;
				}
				ws.close();
				return false;
			};
		});
		</script>
		</head>
		<body>
		<table>
		<tr><td valign="top" width="50%">
		<h1>{{ .Title }}</h1>
		<p>Click "Open" to create a connection to the server, 
		"Send" to send a message to the server and "Close" to close the connection. 
		You can change the message and send multiple times.
		<p>
		<form>
		<button id="open">Open</button>
		<button id="close">Close</button>
		<!--<p><input id="input" type="text" width="400px" value="{{ .DefaultValue }}">
		<button id="send">Send</button>-->
		</form>
		</td><td valign="top" width="50%">
		<div id="output"></div>
		</td></tr></table>
		</body>
		</html>
		`))
