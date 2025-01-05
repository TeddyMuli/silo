import { NextRequest, NextResponse } from 'next/server';
import axios from 'axios';

interface RecaptchaResponseData {
  success: boolean;
  score: number;
  error?: string;
}

interface PostData {
  gRecaptchaToken: string;
}

export async function POST(req: NextRequest) {
  if (req.method !== "POST") {
    return NextResponse.json({ success: false, error: "Method not allowed" }, { status: 405 });
  }

  const secretKey = process.env.RECAPTCHA_SECRET_KEY;
  if (!secretKey) {
    console.error("RECAPTCHA_SECRET_KEY is not set in environment variables.");
    return NextResponse.json({ success: false, error: "Server configuration error" }, { status: 500 });
  }

  let postData: PostData;
  try {
    postData = await req.json() as PostData;
  } catch (error) {
    return NextResponse.json({ success: false, error: "Invalid request body" }, { status: 400 });
  }

  try {
    const response = await axios.post(
      `https://www.google.com/recaptcha/api/siteverify`,
      null,
      {
        params: {
          secret: secretKey,
          response: postData.gRecaptchaToken,
        },
      }
    );

    const data = response.data;

    if (data.success) {
      return NextResponse.json({ success: true, score: data.score }, { status: 200 });
    } else {
      return NextResponse.json({ success: false, error: data['error-codes']?.join(', ') || 'ReCaptcha verification failed' }, { status: 400 });
    }
  } catch (error) {
    console.error("Error verifying ReCaptcha:", error);
    return NextResponse.json({ success: false, error: "ReCaptcha verification failed" }, { status: 500 });
  }
}
