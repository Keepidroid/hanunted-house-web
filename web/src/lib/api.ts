export type CreateJoinResult = {
  roomCode: string;
  playerId: string;
  playerToken: string;
  seat: number;
};

async function request<T>(path: string, init: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: { "Content-Type": "application/json", ...(init.headers || {}) }
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || `HTTP ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export const api = {
  createRoom(nickname: string) {
    return request<CreateJoinResult>("/api/rooms", {
      method: "POST",
      body: JSON.stringify({ nickname })
    });
  },
  joinRoom(roomCode: string, nickname: string) {
    return request<CreateJoinResult>(`/api/rooms/${roomCode}/join`, {
      method: "POST",
      body: JSON.stringify({ nickname })
    });
  },
  startGame(roomCode: string, playerToken: string) {
    return request<{ ok: boolean }>(`/api/rooms/${roomCode}/start`, {
      method: "POST",
      body: JSON.stringify({ playerToken })
    });
  },
  dealDamage(roomCode: string, playerToken: string) {
    return request<{ ok: boolean }>(`/api/rooms/${roomCode}/commands/damage`, {
      method: "POST",
      body: JSON.stringify({
        playerToken,
        kind: "physical",
        allocation: { might: 1 }
      })
    });
  },
  rollTrait(roomCode: string, playerToken: string, trait: "might" | "speed" | "knowledge" | "sanity") {
    return request<{ total: number; dice: number[] }>(`/api/rooms/${roomCode}/commands/roll`, {
      method: "POST",
      body: JSON.stringify({
        playerToken,
        type: "trait",
        trait
      })
    });
  }
};
