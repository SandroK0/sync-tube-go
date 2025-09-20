import { useEffect, useState } from "react";

function App() {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [message, setMessage] = useState<string>("");

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    setWs(socket);

    socket.onopen = () => console.log("Connected to server");
    socket.onmessage = (event) => {
      const msg = event.data;
      setMessage(msg);
    };
    socket.onclose = () => console.log("Disconnected");

    return () => socket.close();
  }, []);

  const sendTest = () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ eventType: "test", data: "Hello from React!" }));
    }
  };

  const sendFoo = () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ eventType: "foo", data: "Hello from React!" }));
    }
  };

  return (
    <div>
      <button onClick={sendTest}>Send Test</button>
      <button onClick={sendFoo}>Send Foo</button>
      <p>Received: {message}</p>
    </div>
  );
}

export default App;
