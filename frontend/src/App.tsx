import axios from "axios";
import { useEffect, useRef, useState } from "react";

type Room = {
  name: string;
};

type Rooms = {
  roomName: Room;
};

type JoinEventData = {
  roomName: string;
  token: string;
};

type MessageEventData = {
  username: string;
  body: string;
};

type RoomReconnectedEventData = {
  roomName: string;
  token: string;
  username: string;
};

type Error = {
  code: string;
  message: string;
};

type ErrorEventData = {
  code: string;
  message: string;
} & Error;

type ServerEvent =
  | { eventType: "room_joined"; data: JoinEventData }
  | { eventType: "message_received"; data: MessageEventData }
  | { eventType: "room_reconnected"; data: RoomReconnectedEventData }
  | { eventType: "error"; data: ErrorEventData };

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
  const [error, setError] = useState<Error | null>(null);

  const usernameInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    setWs(socket);

    socket.onopen = () => {
      console.log("Connected to server");

      const token = localStorage.getItem("token");
      if (token) {
        socket.send(
          JSON.stringify({
            eventType: "reconnect_room",
            data: { token },
          }),
        );
      }
    };

    socket.onmessage = (event: MessageEvent) => {
      try {
        const eventData: ServerEvent = JSON.parse(event.data);
        switch (eventData.eventType) {
          case "room_joined":
            localStorage.setItem("token", eventData.data.token);
            setCurrentRoom(eventData.data.roomName);
            break;
          case "room_reconnected":
            setCurrentRoom(eventData.data.roomName);
            setUsername(eventData.data.username);
            break;
          case "message_received":
            setMessages((prev) => [
              ...prev,
              {
                author: eventData.data.username,
                body: eventData.data.body,
                id: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
              },
            ]);
            break;
          case "error":
            setError(eventData.data);
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
          eventType: "join_room",
          data: {
            roomName,
            username: username,
          },
        }),
      );
    }
  };

  const createRoom = async () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          eventType: "create_room",
          data: {
            roomName: createRoomName,
            username: username,
          },
        }),
      );
    }

    fetchRooms();
  };
  const leaveRoom = async () => {
    const token = localStorage.getItem("token");

    if (!token) {
      console.error("no token");
      return;
    }

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          eventType: "leave_room",
          data: {
            roomName: currentRoom,
            username: username,
            token,
          },
        }),
      );
    }

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
          eventType: "send_message",
          data: {
            roomName: currentRoom,
            username: username,
            body: messageBody,
          },
        }),
      );
    }
  };

  return (
    <div className="flex gap-10 flex-col items-center justify-center relative w-screen min-h-screen ">
      {error && (
        <div className="text-red-800 absolute bg-amber-50 top-10">
          {error.code}:{error.message}
        </div>
      )}
      <div className="flex items-center gap-5">
        <input
          placeholder="set username"
          className="focus:ring-0 outline-0 border rounded p-2"
          ref={usernameInputRef}
        ></input>
        <button
          disabled={username !== ""}
          onClick={() => {
            if (usernameInputRef.current) {
              setUsername(usernameInputRef.current.value);
            }
          }}
          className="disabled:opacity-50"
        >
          Set
        </button>
      </div>
      <div className="mt-40 flex gap-10 ">
        <div>
          <div className="flex justify-between">
            <input
              type="text"
              className="focus:ring-0 outline-0 border rounded p-2"
              onChange={(e) => setCreateRoomName(e.target.value)}
              value={createRoomName}
            />
            <button onClick={createRoom}>Create Room</button>
            <button onClick={leaveRoom}>Leave Room</button>
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
              <p key={message.id}>
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
