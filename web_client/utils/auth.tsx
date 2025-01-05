import { cookies } from 'next/headers';
import jwt from 'jsonwebtoken';
import { redirect } from 'next/navigation';

export default function ProtectedPage() {
  const cookieStore = cookies();
  const token = cookieStore.get('auth_token')?.value;

  if (!token) {
    redirect("/login")
  }

  try {
    const user = jwt.verify(token, process.env.JWT_SECRET!);
  } catch (err) {
    console.error("Error getting token: ", err)
  }
}
