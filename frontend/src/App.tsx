import axios from "axios";
import "./App.css";
import brick from "./assets/brick.png";
import background from "./assets/background.jpg";
import character from "./assets/character.png"
import "xp.css/dist/XP.css";

import { Stage, Sprite } from "@pixi/react";
import { useEffect, useState } from "react";

const App = () => {
  let [world, setWorld] = useState<String[][]>([]);
  let [characterPos, setCharacterPos] = useState<number[]>([5, 5]);

  const tileSize = 16;

  useEffect(() => {
    loadWorld();
  }, []);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      switch (event.key) {
        case "w":
          if (characterPos[1] == 0) {
            setCharacterPos([characterPos[0], world.length - 1])
            loadWorld();
          } else if (worldAt(characterPos[0], characterPos[1] - 1) === "STONE_FLOOR") {
            setCharacterPos([characterPos[0], characterPos[1] - 1]);
          }
          break;
        case "a":
          if (characterPos[0] == 0) {
            setCharacterPos([world[0].length - 1, characterPos[1]])
            loadWorld();
          } else if (worldAt(characterPos[0] - 1, characterPos[1]) === "STONE_FLOOR") {
            setCharacterPos([characterPos[0] - 1, characterPos[1]]);
          }
          break;
        case "s":
          if (characterPos[1] == world.length - 1) {
            setCharacterPos([characterPos[0], 0])
            loadWorld();
          } else if (worldAt(characterPos[0], characterPos[1] + 1) === "STONE_FLOOR") {
            setCharacterPos([characterPos[0], characterPos[1] + 1]);
          }
          break;
        case "d":
          if (characterPos[0] == world[0].length - 1) {
            setCharacterPos([0, characterPos[1]])
            loadWorld();
          } else if (worldAt(characterPos[0] + 1, characterPos[1]) === "STONE_FLOOR") {
            setCharacterPos([characterPos[0] + 1, characterPos[1]]);
          }
          break;
        case "f":
          console.log(characterPos)
          break;
      }
    };

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [characterPos, world]);

  const loadWorld = () => {
    axios
      .get("http://localhost:8080/api/chunk", {
        params: {
          x: 1,
          y: 2,
        },
      })
      .then((response) => {
        setWorld(response.data.rows);
      });
  };

  const worldAt = (x: number, y: number) => {
    if (x < 0 || y < 0 || y > world.length || x > world[0].length) {
      return ""
    }
    return world[y][x]
  }

  return (
    <div className="window">
      <div className="title-bar">
        <div className="title-bar-text">Dungeon Viewer</div>
        <div className="title-bar-controls">
          <button aria-label="Minimize" />
          <button aria-label="Maximize" />
          <button aria-label="Close" />
        </div>
      </div>

      <div className="window-body">
        <div className="field-row" style={{ justifyContent: "center" }}>
          <Stage width={tileSize * 48} height={tileSize * 36} options={{ background: 0x1099bb }}>
            <Sprite image={background} x={-600} y={-200} />
            {world.flatMap((row, y) =>
              row.flatMap((cell, x) => {
                const spriteX = x * tileSize;
                const spriteY = y * tileSize;
                if (cell === "STONE_FLOOR" && spriteX < 800 && spriteY < 600) {
                  return (
                    <Sprite
                      image={brick}
                      x={x * tileSize}
                      y={y * tileSize}
                    />
                  );
                }
              }),
            )}
            <Sprite image={character} x={characterPos[0] * tileSize} y={characterPos[1] * tileSize} />
          </Stage>
        </div>
      </div>
    </div>
  );
};

export default App;
