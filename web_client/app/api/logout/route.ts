import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  const response = NextResponse.redirect(new URL('/login', request.url));
  
  // Clear the auth_token cookie
  response.cookies.set('auth_token', '', { maxAge: -1, path: '/' });

  return response;
}
