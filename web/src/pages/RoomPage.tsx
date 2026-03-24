import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { ConnectionBadge } from "../components/ConnectionBadge";
import { api } from "../lib/api";
import { useRoomStore } from "../store/room";
import { useSessionStore } from "../store/session";

export function RoomPage() {
  const { roomCode = "" } = useParams();
  const nav = useNavigate();
  const session = useSessionStore();
  const latest = useRoomStore((s) => s.latest);
  const setConnected = useRoomStore((s) => s.setConnected);
  const setSnapshot = useRoomStore((s) => s.setSnapshot);
  const [error, setError] = useState("");

  const isHost = useMemo(() => latest?.players.find((p) => p.isCurrent)?.isHost ?? false, [latest]);

  useEffect(() => {
    if (!session.playerId || session.roomCode !== roomCode.toUpperCase()) {
      nav("/");
      return;
    }

    const ws = new WebSocket(`${location.protocol === "https:" ? "wss" : "ws"}://${location.host}/ws/${roomCode}?playerId=${session.playerId}`);
    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onerror = () => setConnected(false);
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === "snapshot") {
        setSnapshot(data.payload);
      }
    };

    return () => {
      ws.close();
    };
  }, [nav, roomCode, session.playerId, session.roomCode, setConnected, setSnapshot]);

  async function onStart() {
    setError("");
    try {
      await api.startGame(roomCode, session.playerToken);
    } catch (err) {
      setError(err instanceof Error ? err.message : "开始失败");
    }
  }

  return (
    <main className="container">
      <header className="room-header">
        <h2>房间 {roomCode.toUpperCase()}</h2>
        <ConnectionBadge />
      </header>
      <section className="panel">
        <p>阶段：{latest?.phase ?? "lobby"}</p>
        <p>当前回合座位：{latest?.turnSeat ?? 0}</p>
        <ul>
          {latest?.players.map((p) => (
            <li key={p.playerId}>
              Seat {p.seat} - {p.nickname}
              {p.isHost ? "（房主）" : ""}
              {p.isCurrent ? "（你）" : ""}
              {p.connected ? " ✅" : " ❌"}
            </li>
          ))}
        </ul>
        {isHost && latest?.phase === "lobby" && <button onClick={onStart}>开始游戏（至少3人）</button>}
      </section>
      {error && <p className="error">{error}</p>}
    </main>
  );
}
