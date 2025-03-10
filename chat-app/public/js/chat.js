document.addEventListener("DOMContentLoaded", function () {
    const socket = new WebSocket("ws://localhost:8080/ws");
    const chatBox = document.getElementById("chat-box");
    const messageInput = document.getElementById("message");
    const sendButton = document.getElementById("send");

    // WebSocket が開いたとき
    socket.onopen = function () {
        console.log("WebSocket 接続が確立しました。");
    };

    // メッセージを受信したとき
    socket.onmessage = function (event) {
        const message = document.createElement("div");
        message.textContent = event.data;
        chatBox.appendChild(message);
        chatBox.scrollTop = chatBox.scrollHeight;
    };

    // WebSocket が閉じたとき
    socket.onclose = function () {
        console.log("WebSocket 接続が閉じられました。");
    };

    // 送信ボタンをクリックしたとき
    sendButton.addEventListener("click", function () {
        if (messageInput.value.trim() !== "") {
            socket.send(messageInput.value);
            messageInput.value = "";
        }
    });

    // Enterキーで送信できるように
    messageInput.addEventListener("keypress", function (event) {
        if (event.key === "Enter") {
            sendButton.click();
        }
    });
});
