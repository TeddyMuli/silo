"use client";

import { QueryClientProvider, QueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { GoogleReCaptchaProvider } from "react-google-recaptcha-v3";

const ReactQueryProvider = ({ children }: { children: React.ReactNode }) => {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <GoogleReCaptchaProvider reCaptchaKey={`${process.env.NEXT_PUBLIC_RECAPTCHA_SITE_KEY}`}>
        {children}
      </GoogleReCaptchaProvider>
    </QueryClientProvider>
  );
};

export default ReactQueryProvider;
