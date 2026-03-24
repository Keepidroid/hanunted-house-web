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

export type ExplorerTrait = {
  trait: "might" | "speed" | "knowledge" | "sanity";
  trackIndex: number;
  value: number;
  critical: boolean;
};

export type ExplorerView = {
  explorerId: string;
  character: string;
  side: "heroes" | "unknown";
  dead: boolean;
  traits: ExplorerTrait[];
};

export type PlayerView = {
  roomCode: string;
  phase: "lobby" | "preHaunt" | "postHaunt" | "finished";
  players: Seat[];
  turnSeat: number;
  version: number;
  explorer?: ExplorerView;
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
