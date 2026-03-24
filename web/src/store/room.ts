import { create } from "zustand";

export type Seat = {
  seat: number;
  playerId: string;
  nickname: string;
  connected: boolean;
  joinedAt: string;
  isHost: boolean;
  isCurrent: boolean;
};

export type PlayerView = {
  roomCode: string;
  phase: "lobby" | "preHaunt" | "finished";
  players: Seat[];
  turnSeat: number;
  version: number;
};

type RoomState = {
  connected: boolean;
  latest?: PlayerView;
  setConnected: (v: boolean) => void;
  setSnapshot: (v: PlayerView) => void;
};

export const useRoomStore = create<RoomState>((set) => ({
  connected: false,
  setConnected(v) {
    set({ connected: v });
  },
  setSnapshot(v) {
    set({ latest: v });
  }
}));
