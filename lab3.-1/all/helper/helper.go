package main

import (
	"html/template"
	"net/http"

	log "github.com/mgutz/logxi/v1"
)

const INDEX_HTML = `
<html>
<head>
<title>Летучка</title>
<meta charset="utf8">
<script>

let socket1 = new WebSocket("ws://46.138.248.69:9457");
let socket2 = new WebSocket("ws://46.138.248.69:9460");
let socket3 = new WebSocket("ws://185.102.139.168:9457");
let socket4 = new WebSocket("ws://185.102.139.169:9457");

socket1.onmessage = function(event) {
  //alert("[message] Данные получены с сервера: ${event.data}");
  document.getElementById("s1").innerHTML = event.data; 
};

socket2.onmessage = function(event) {
  //alert("[message] Данные получены с сервера: ${event.data}");
  document.getElementById("s2").innerHTML = event.data; 
};

socket3.onmessage = function(event) {
  //alert("[message] Данные получены с сервера: ${event.data}");
  document.getElementById("s3").innerHTML = event.data; 
};

socket4.onmessage = function(event) {
  //alert("[message] Данные получены с сервера: ${event.data}");
  document.getElementById("s4").innerHTML = event.data; 
};


socket1.onerror = function(error) {
  document.getElementById("s1").innerHTML = "disconnected";
};

socket2.onerror = function(error) {
  document.getElementById("s2").innerHTML = "disconnected";
};

socket3.onerror = function(error) {
  document.getElementById("s3").innerHTML = "disconnected";
};

socket4.onerror = function(error) {
  document.getElementById("s4").innerHTML = "disconnected";
};

</script>
</head>
<body>
<div id="s1">Limonov</div><br>
<div id="s2">Limonov</div><br>
<div id="s3">Limonov</div><br>
<div id="s4">Limonov</div><br>
<input name="usersText" id="usersText" placeholder="Type text for 1st peer">
<button type="button" id="post-btn" onclick="fetching()">Send</button>
<script>
  function fetching() {
    let text = document.getElementsByName("usersText")[0].value
    let data = {
      tkey: text
    }
    fetch("http://46.138.248.69:9455/h", {
	  mode: 'no-cors',
      method: "POST",
      body: JSON.stringify(data)
    })
  }
</script>
</body>
</html>
    `

var indexHtml = template.Must(template.New("index").Parse(INDEX_HTML))

func serveClient(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	log.Info("got request", "Method", request.Method, "Path", path)
	if path != "/" && path != "/index.html" {
		log.Error("invalid path", "Path", path)
		response.WriteHeader(http.StatusNotFound)
	} else if err := indexHtml.Execute(response, nil); err != nil {
		log.Error("HTML creation failed", "error", err)
	} else {
		log.Info("response sent to client successfully")
	}
}

func main() {
	http.HandleFunc("/", serveClient)
	log.Info("starting listener")
	log.Error("listener failed", "error", http.ListenAndServe("192.168.31.116:8122", nil))
}
