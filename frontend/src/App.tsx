import axios from "axios";
import "./App.css";
import brick from "./assets/brick.png";
import background from "./assets/background.jpg";
import character from "./assets/character.png";
import character2 from "./assets/character2.png";
import "xp.css/dist/XP.css";

import { Stage, Sprite } from "@pixi/react";
import { useEffect, useRef, useState } from "react";
import { incomingMessage } from "./models/message";

const App = () => {
  const [world, setWorld] = useState<string[][]>([]);
  const [worldPos, setWorldPos] = useState<number[]>([0, 0]);
  const [viewOffset, setViewOffset] = useState<number[]>([0, 0]);
  const [characterPos, setCharacterPos] = useState<number[]>([5, 5]);
  const [zoom, setZoom] = useState(2);
  const [userId, setUserId] = useState("")
  const [userToPosition, setUserToPosition] = useState<Map<string, number[]>>(new Map()) 
  const [ws, setWs] = useState<WebSocket | undefined>(undefined);

  const userIdRef = useRef(userId);

  const BASE_URL = "http://localhost:8080/api";
  const CHUNK_URL = `${BASE_URL}/chunk`;
  const WS_URL = `${BASE_URL}/ws`;

  const WIDTH = 48;
  const HEIGHT = 36;
  const TILESIZE = 16;

  useEffect(() => {
    const socket = new WebSocket(WS_URL);

    socket.onopen = () => {
      console.log("Connected to WebSocket server");
      setWs(socket);
    };

    socket.onmessage = (event) => {
      const message = JSON.parse(event.data) as incomingMessage;
      switch (message.type) {
        case "userId":
          setUserId(message.userId)
          break;
        case "position":
          if (message.userId === userIdRef.current) {
            setCharacterPosWithView(message.data.x, message.data.y);
          } else {
            setUserToPosition((prev) => {
              const newMap = new Map(prev);
              newMap.set(message.userId, [message.data.x, message.data.y]);
              return newMap;
            });
          }
          break;
        case "playerExit":
          setUserToPosition((prev) => {
            const newMap = new Map(prev);
            newMap.delete(message.userId)
            return newMap
          });
          break;
      }
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    socket.onclose = () => {
      console.log("WebSocket connection closed");
    };
  }, []);

  useEffect(() => {
    loadWorld(worldPos[0], worldPos[1]);
  }, [ws]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      switch (event.key) {
        case "w":
          if (characterPos[1] == 0) {
            loadWorld(worldPos[0], worldPos[1] - 1);
            sendPositionMessage(characterPos[0], world.length - 1);
          } else if (
            worldAt(characterPos[0], characterPos[1] - 1) === "STONE_FLOOR"
          ) {
            sendPositionMessage(characterPos[0], characterPos[1] - 1);
          }
          break;
        case "a":
          if (characterPos[0] == 0) {
            loadWorld(worldPos[0] - 1, worldPos[1]);
            sendPositionMessage(world[0].length - 1, characterPos[1]);
          } else if (
            worldAt(characterPos[0] - 1, characterPos[1]) === "STONE_FLOOR"
          ) {
            sendPositionMessage(characterPos[0] - 1, characterPos[1]);
          }
          break;
        case "s":
          if (characterPos[1] == world.length - 1) {
            loadWorld(worldPos[0], worldPos[1] + 1);
            sendPositionMessage(characterPos[0], 0);
          } else if (
            worldAt(characterPos[0], characterPos[1] + 1) === "STONE_FLOOR"
          ) {
            sendPositionMessage(characterPos[0], characterPos[1] + 1);
          }
          break;
        case "d":
          if (characterPos[0] == world[0].length - 1) {
            loadWorld(worldPos[0] + 1, worldPos[1]);
            sendPositionMessage(0, characterPos[1]);
          } else if (
            worldAt(characterPos[0] + 1, characterPos[1]) === "STONE_FLOOR"
          ) {
            sendPositionMessage(characterPos[0] + 1, characterPos[1]);
          }
          break;
        case "z":
          if (zoom === 1) {
            setZoom(2);
          } else {
            setZoom(1);
          }
          break;
      }
    };

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [characterPos, world, zoom]);

  useEffect(() => {
    userIdRef.current = userId;
  }, [userId]);

  const loadWorld = (x: number, y: number) => {
    if (ws == null) {
      return;
    }

    ws.send(
      JSON.stringify({
        type: "subscribe",
        data: {
          topic: `${x},${y}`,
        },
      }),
    );

    setUserToPosition(new Map())

    axios
      .get(CHUNK_URL, {
        params: {
          x,
          y,
        },
      })
      .then((response) => {
        setWorld(response.data.world);
        setWorldPos([x, y]);
      });
  };

  const worldAt = (x: number, y: number) => {
    if (x < 0 || y < 0 || y > world.length || x > world[0].length) {
      return "";
    }
    return world[y][x];
  };

  const setCharacterPosWithView = (x: number, y: number) => {
    // TODO: Fix this crap up
    setCharacterPos([x, y]);

    let viewX = Math.min(
      Math.max(x - WIDTH / zoom / 2, 0),
      WIDTH - WIDTH / zoom,
    );
    let viewY = Math.min(
      Math.max(y - HEIGHT / zoom / 2, 0),
      HEIGHT - HEIGHT / zoom,
    );

    setViewOffset([viewX, viewY]);
  };

  const sendPositionMessage = (x: number, y: number) => {
    ws?.send(
      JSON.stringify({
        type: "position",
        data: {
          x,
          y,
        },
      }),
    );
  };

  return (
    <div className="window">
      <div className="title-bar">
        <div className="title-bar-text">
          World Dungeon ({worldPos[0]}, {worldPos[1]})
        </div>
        <div className="title-bar-controls">
          <button aria-label="Minimize" />
          <button aria-label="Maximize" />
          <button aria-label="Close" />
        </div>
      </div>

      <div className="window-body">
        <div className="field-row" style={{ justifyContent: "center" }}>
          <Stage
            width={TILESIZE * WIDTH}
            height={TILESIZE * HEIGHT}
            options={{ background: 0x1099bb }}
          >
            <Sprite image={background} x={-600} y={-200} />
            {world.flatMap((row, y) =>
              row.flatMap((cell, x) => {
                const spriteX = x * TILESIZE * zoom;
                const spriteY = y * TILESIZE * zoom;
                if (cell === "STONE_FLOOR") {
                  return (
                    <Sprite
                      image={brick}
                      x={spriteX - viewOffset[0] * TILESIZE * zoom}
                      y={spriteY - viewOffset[1] * TILESIZE * zoom}
                      width={TILESIZE * zoom}
                      height={TILESIZE * zoom}
                    />
                  );
                }
              }),
            )}
            <Sprite
              image={character}
              x={
                characterPos[0] * TILESIZE * zoom -
                viewOffset[0] * TILESIZE * zoom
              }
              y={
                characterPos[1] * TILESIZE * zoom -
                viewOffset[1] * TILESIZE * zoom
              }
              width={TILESIZE * zoom}
              height={TILESIZE * zoom}
            />
            {Array.from(userToPosition.entries()).map(([_, pos]) => {
              return <Sprite
              image={character2}
              x={
                pos[0] * TILESIZE * zoom -
                viewOffset[0] * TILESIZE * zoom
              }
              y={
                pos[1] * TILESIZE * zoom -
                viewOffset[1] * TILESIZE * zoom
              }
              width={TILESIZE * zoom}
              height={TILESIZE * zoom}
            />
            })}
          </Stage>
        </div>
      </div>
    </div>
  );
};

export default App;
