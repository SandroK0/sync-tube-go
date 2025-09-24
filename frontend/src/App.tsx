import axios from "axios";
import { useEffect, useState } from "react";

type Room = {
  name: string;
};

type Rooms = {
  roomName: Room;
};

type EventType = "join" | "message";

type JoinEventData = {
  roomName: string;
};

type MessageEventData = {
  username: string;
  body: string;
};

type Event<T extends EventType = EventType> = T extends "join"
  ? { eventType: "join"; data: JoinEventData }
  : { eventType: "message"; data: MessageEventData };

type Message = {
  author: string;
  body: string;
  id: string;
};

function App() {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [rooms, setRooms] = useState<Rooms | null>(null);
  const [createRoomName, setCreateRoomName] = useState<string>("");
  const [currentRoom, setCurrentRoom] = useState<string>("");
  const [username, setUsername] = useState<string>("");

  useEffect(() => {
    const randomString = Math.random().toString(36).substring(2, 10);
    setUsername(randomString);

    const socket = new WebSocket("ws://localhost:8080/ws");
    setWs(socket);

    socket.onopen = () => console.log("Connected to server");

    socket.onmessage = (event: MessageEvent) => {
      try {
        const eventData: Event = JSON.parse(event.data);
        switch (eventData.eventType) {
          case "join":
            setCurrentRoom(eventData.data.roomName); 
            break;
          case "message":
            setMessages((prev) => [
              ...prev,
              {
                author: eventData.data.username,
                body: eventData.data.body,
                id: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
              },
            ]);
            break;
          default:
        }
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    fetchRooms();

    socket.onclose = () => console.log("Disconnected");
    return () => socket.close();
  }, []);

  const joinRoom = async (roomName: string) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          eventType: "join",
          data: {
            roomName,
            username: username,
          },
        })
      );
    }
  };

  const createRoom = async () => {
    await axios.post("http://localhost:8080/rooms/create", {
      roomName: createRoomName,
    });

    fetchRooms();
  };

  const fetchRooms = async () => {
    const response = await axios.get("http://localhost:8080/rooms");
    setRooms(response.data);
  };

  const sendMessage = (messageBody: string) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          eventType: "message",
          data: {
            roomName: currentRoom,
            username: username,
            body: messageBody,
          },
        })
      );
    }
  };

  return (
    <div className="flex gap-10  justify-center w-screen min-h-screen ">
      <div className="mt-40 flex gap-10">
        <div>
          <div className="flex justify-between">
            <input
              type="text"
              className="focus:ring-0 outline-0 border rounded p-2"
              onChange={(e) => setCreateRoomName(e.target.value)}
              value={createRoomName}
            />
            <button onClick={createRoom}>Create Room</button>
          </div>
          <div className="flex flex-col gap-2 border rounded w-100 h-100 mt-5">
            {rooms &&
              Object.keys(rooms).map((key) => (
                <div
                  key={key}
                  className="flex justify-between border-b px-4 py-1 items-center"
                >
                  <div>{key}</div>
                  <button onClick={() => joinRoom(key)}>Join</button>
                </div>
              ))}
          </div>
          <div>Current Room:{currentRoom}</div>
          <div>Username:{username}</div>
        </div>
        <div className="flex flex-col">
          <h1 className="p-2 font-medium">Messages</h1>
          <div className="flex flex-col gap-5 border rounded w-100 h-100 p-5 mt-5">
            {messages.map((message: Message) => (
              <p>
                {message.author}:{message.body}
              </p>
            ))}
          </div>
          <div className="w-full flex justify-between gap-3 mt-2">
            <input
              type="text"
              className="focus:ring-0 outline-0 border rounded p-2 w-full"
              onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => {
                if (e.key === "Enter") {
                  const value = (e.target as HTMLInputElement).value;
                  if (value.trim() === "") {
                    return;
                  }
                  sendMessage(value);
                  (e.target as HTMLInputElement).value = "";
                }
              }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
