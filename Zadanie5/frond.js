document.addEventListener('DOMContentLoaded', function(){
    document.getElementById("btn").addEventListener("click", () => {
        let limit = document.getElementById("limit").value;
        let flow = document.getElementById("flow").value;
        let count = document.getElementById("count").value;
        document.getElementById("result").innerHTML = "";

        let url = `ws://192.168.81.53:8888/ws?limit=${limit}&flow=${flow}&count=${count}`;
        let socket = new WebSocket(url);
        socket.onopen = () => {
            socket.send(url)
            console.log("Socket open connection!");
        };
        socket.onmessage = event => {
            document.getElementById('result').innerHTML += event.data+"  ";
        }
        socket.onclose = event => {
            console.log("Socket closed connection: ", event);
        };
        socket.onerror = error => {
            console.log("Socket error: ", error);
        };
    });
}, false);
