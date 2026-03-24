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
  const [rollResult, setRollResult] = useState<string>("");

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

  async function onDamage() {
    setError("");
    try {
      await api.dealDamage(roomCode, session.playerToken);
    } catch (err) {
      setError(err instanceof Error ? err.message : "伤害结算失败");
    }
  }

  async function onRollMight() {
    setError("");
    try {
      const res = await api.rollTrait(roomCode, session.playerToken, "might");
      setRollResult(`Might 掷骰 => [${res.dice.join(", ")}] 总计 ${res.total}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "掷骰失败");
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

      {latest?.explorer && (
        <section className="panel">
          <h3>我的探索者</h3>
          <p>角色：{latest.explorer.character}</p>
          <p>阵营：{latest.explorer.side}</p>
          <ul>
            {latest.explorer.traits.map((t) => (
              <li key={t.trait}>
                {t.trait}: 值 {t.value}（轨道位置 {t.trackIndex}）{t.critical ? " [critical]" : ""}
              </li>
            ))}
          </ul>
          {latest.phase !== "lobby" && (
            <div className="actions" style={{ gridTemplateColumns: "1fr 1fr" }}>
              <button onClick={onDamage}>测试：受到 1 点 physical 伤害</button>
              <button onClick={onRollMight}>测试：Might Trait Roll</button>
            </div>
          )}
          {rollResult && <p>{rollResult}</p>}
        </section>
      )}

      {error && <p className="error">{error}</p>}
    </main>
  );
}
