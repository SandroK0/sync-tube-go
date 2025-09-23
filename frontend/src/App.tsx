import axios from "axios";
import { useEffect, useState } from "react";

type Room = {
  name: string;
};

type Rooms = {
  roomName: Room;
};

function App() {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [message, setMessage] = useState<string>("");
  const [rooms, setRooms] = useState<Rooms | null>(null);
  const [createRoomName, setCreateRoomName] = useState<string>("");

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    setWs(socket);

    socket.onopen = () => console.log("Connected to server");
    socket.onmessage = (event) => {
      const msg = event.data;
      setMessage(msg);
    };
    socket.onclose = () => console.log("Disconnected");
    fetchRooms();

    return () => socket.close();
  }, []);

  const joinRoom = async (roomName: string) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          eventType: "join",
          data: {
            roomName,
            username: "Sandro",
          },
        }),
      );
    }
  };

  const createRoom = async () => {
    const response = await axios.post("http://localhost:8080/rooms/create", {
      roomName: createRoomName,
    });

    fetchRooms();
    console.log(response.data);
  };

  const fetchRooms = async () => {
    const response = await axios.get("http://localhost:8080/rooms");
    console.log(response.data);
    setRooms(response.data);
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
          <div>Current Name:</div>
        </div>
        <div className="flex flex-col">
          <h1 className="p-2 font-medium">Messages</h1>
          <div className="flex flex-col gap-5 border rounded w-100 h-100 p-5 mt-5">
            <p>{message}</p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
