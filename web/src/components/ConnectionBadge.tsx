import { useRoomStore } from "../store/room";

export function ConnectionBadge() {
  const connected = useRoomStore((s) => s.connected);
  return <span className={connected ? "badge ok" : "badge err"}>{connected ? "已连接" : "连接中断"}</span>;
}
