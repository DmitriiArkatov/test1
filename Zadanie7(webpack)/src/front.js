document.addEventListener('DOMContentLoaded', function(){
    document.getElementById("btn").addEventListener("click", () => {
        // привязываем переменные к элементам DOM
        let limit = document.getElementById("limit").value;
        let flow = document.getElementById("flow").value;
        let count = document.getElementById("count").value;
        document.getElementById("result").innerHTML = "";
        // Извлекаем базовую часть URL (хост и порт)
        let host = window.location.host;
        // Создаем URL, объединяя базовую часть и относительный путь
        let url = `ws://${host}/ws?limit=${limit}&flow=${flow}&count=${count}`;
        let socket = new WebSocket(url);
        //открываем соединение
        socket.onopen = () => {
            socket.send(url)
            console.log("Socket open connection!");
        };
        //получаем ответ от сервера и выводим его на экран
        socket.onmessage = event => {
            document.getElementById('result').innerHTML += event.data+"  ";
        }
        //закрываем соединение
        socket.onclose = event => {
            console.log("Socket closed connection: ", event);
        };
        //ошибка соединения
        socket.onerror = error => {
            console.log("Socket error: ", error);
        };
    });
}, false);