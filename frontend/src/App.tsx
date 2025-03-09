import axios from "axios";
import "./App.css";
import brick from "./assets/brick.png";
import background from "./assets/background.jpg";
import character from "./assets/character.png";
import "xp.css/dist/XP.css";

import { Stage, Sprite } from "@pixi/react";
import { useEffect, useState } from "react";

const App = () => {
  let [world, setWorld] = useState<String[][]>([]);
  let [worldPos, setWorldPos] = useState<number[]>([0, 0]);
  let [viewOffset, setViewOffset] = useState<number[]>([0, 0]);
  let [characterPos, setCharacterPos] = useState<number[]>([5, 5]);
  let [zoom, setZoom] = useState(2)

  const WIDTH = 48;
  const HEIGHT = 36;
  const TILESIZE = 16;

  useEffect(() => {
    loadWorld(worldPos[0], worldPos[1]);
  }, []);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      switch (event.key) {
        case "w":
          if (characterPos[1] == 0) {
            setCharacterPosWithView(characterPos[0], world.length - 1);
            loadWorld(worldPos[0], worldPos[1] - 1);
          } else if (
            true || worldAt(characterPos[0], characterPos[1] - 1) === "STONE_FLOOR"
          ) {
            setCharacterPosWithView(characterPos[0], characterPos[1] - 1);
          }
          break;
        case "a":
          if (characterPos[0] == 0) {
            setCharacterPosWithView(world[0].length - 1, characterPos[1]);
            loadWorld(worldPos[0] - 1, worldPos[1]);
          } else if (
            true || worldAt(characterPos[0] - 1, characterPos[1]) === "STONE_FLOOR"
          ) {
            setCharacterPosWithView(characterPos[0] - 1, characterPos[1]);
          }
          break;
        case "s":
          if (characterPos[1] == world.length - 1) {
            setCharacterPosWithView(characterPos[0], 0);
            loadWorld(worldPos[0], worldPos[1] + 1);
          } else if (
            true || worldAt(characterPos[0], characterPos[1] + 1) === "STONE_FLOOR"
          ) {
            setCharacterPosWithView(characterPos[0], characterPos[1] + 1);
          }
          break;
        case "d":
          if (characterPos[0] == world[0].length - 1) {
            setCharacterPosWithView(0, characterPos[1]);
            loadWorld(worldPos[0] + 1, worldPos[1]);
          } else if (
            true || worldAt(characterPos[0] + 1, characterPos[1]) === "STONE_FLOOR"
          ) {
            setCharacterPosWithView(characterPos[0] + 1, characterPos[1]);
          }
          break;
        case "z":
          if (zoom === 1) {
            setZoom(2)
          } else {
            setZoom(1)
          }
          break;
      }
    };

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [characterPos, world, zoom]);

  const loadWorld = (x: number, y: number) => {
    axios
      .get("http://localhost:8080/api/chunk", {
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
    setCharacterPos([x, y]);

    let viewX = Math.min(Math.max(x - (WIDTH / zoom / 2), 0), WIDTH - (WIDTH / zoom));
    console.log(x - (WIDTH / zoom / 2))
    let viewY = Math.min(Math.max(y - (HEIGHT / zoom / 2), 0), HEIGHT - (HEIGHT / zoom));

    setViewOffset([viewX, viewY]);
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
              x={characterPos[0] * TILESIZE * zoom - viewOffset[0] * TILESIZE * zoom}
              y={characterPos[1] * TILESIZE * zoom - viewOffset[1] * TILESIZE * zoom}
              width={TILESIZE * zoom}
              height={TILESIZE * zoom}
            />
          </Stage>
        </div>
      </div>
    </div>
  );
};

export default App;
