import { create } from "zustand";

const KEY = "hanunted-session";

type Session = {
  nickname: string;
  roomCode: string;
  playerId: string;
  playerToken: string;
  setSession: (data: Omit<Session, "setSession" | "clear">) => void;
  clear: () => void;
};

function load() {
  try {
    const raw = localStorage.getItem(KEY);
    if (!raw) return null;
    return JSON.parse(raw) as Omit<Session, "setSession" | "clear">;
  } catch {
    return null;
  }
}

const initial = load() || { nickname: "", roomCode: "", playerId: "", playerToken: "" };

export const useSessionStore = create<Session>((set) => ({
  ...initial,
  setSession(data) {
    localStorage.setItem(KEY, JSON.stringify(data));
    set(data);
  },
  clear() {
    localStorage.removeItem(KEY);
    set({ nickname: "", roomCode: "", playerId: "", playerToken: "" });
  }
}));
