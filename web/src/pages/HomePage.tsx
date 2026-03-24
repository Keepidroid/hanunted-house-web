import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { api } from "../lib/api";
import { useSessionStore } from "../store/session";

export function HomePage() {
  const nav = useNavigate();
  const setSession = useSessionStore((s) => s.setSession);
  const [nickname, setNickname] = useState("");
  const [roomCode, setRoomCode] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function onCreate(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const out = await api.createRoom(nickname);
      setSession({ nickname, ...out });
      nav(`/room/${out.roomCode}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "创建失败");
    } finally {
      setLoading(false);
    }
  }

  async function onJoin(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const out = await api.joinRoom(roomCode.trim().toUpperCase(), nickname);
      setSession({ nickname, ...out });
      nav(`/room/${out.roomCode}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "加入失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="container">
      <h1>山屋惊魂（第三版）在线版</h1>
      <p>输入昵称后创建房间或输入房间码加入。</p>
      <div className="panel">
        <label>昵称</label>
        <input value={nickname} onChange={(e) => setNickname(e.target.value)} placeholder="例如：夜行者" />
      </div>
      <div className="actions">
        <form onSubmit={onCreate}>
          <button disabled={loading}>创建房间</button>
        </form>
        <form onSubmit={onJoin} className="join-form">
          <input value={roomCode} onChange={(e) => setRoomCode(e.target.value)} placeholder="房间码" />
          <button disabled={loading}>加入房间</button>
        </form>
      </div>
      {error && <p className="error">{error}</p>}
    </main>
  );
}
