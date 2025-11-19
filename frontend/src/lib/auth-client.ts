import { createAuthClient } from "better-auth/client";
import { magicLinkClient, jwtClient } from "better-auth/client/plugins";

export const authClient = createAuthClient({
	baseURL: "http://localhost:4321",
	plugins: [
		magicLinkClient(),
		jwtClient(),
  ]
});
