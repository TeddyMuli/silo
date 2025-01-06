"use client";

import { useUser } from "@/context/UserContext";

export default function Home() {
  const { user } = useUser()

  return (
    <div className="container">
      <h1>Silo</h1>
      <p>User: {user?.user_id}</p>
    </div>
  );
}
