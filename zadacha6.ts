document.addEventListener('DOMContentLoaded', () => {
    document.getElementById("btn")!.addEventListener("click", () => {
        // привязываем переменные к элементам DOM
        const limitInput:HTMLInputElement = document.getElementById("limit") as HTMLInputElement;
        const flowInput:HTMLInputElement = document.getElementById("flow") as HTMLInputElement;
        const countInput:HTMLInputElement = document.getElementById("count") as HTMLInputElement;
        document.getElementById("result")!.innerHTML = "";

        // Извлекаем базовую часть URL (хост и порт)
        const host:string = window.location.host;

        // Создаем URL, объединяя базовую часть и относительный путь
        const limit:string = limitInput.value;
        const flow:string = flowInput.value;
        const count:string = countInput.value;
        const url:string = `ws://${host}/ws?limit=${limit}&flow=${flow}&count=${count}`;

        const socket:WebSocket = new WebSocket(url);

        //открываем соединение
        socket.onopen = () => {
            socket.send(url);
            console.log("Socket open connection!");
        };

        //получаем ответ от сервера и выводим его на экран
        socket.onmessage = (event: MessageEvent) => {
            document.getElementById('result')!.innerHTML += event.data + "  ";
        };

        //закрываем соединение
        socket.onclose = (event: CloseEvent) => {
            console.log("Socket closed connection: ", event);
        };

        //ошибка соединения
        socket.onerror = (error: Event) => {
            console.log("Socket error: ", error);
        };
    });
}, false);
