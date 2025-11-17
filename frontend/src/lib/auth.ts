import { betterAuth } from "better-auth";
import { mongodbAdapter } from "better-auth/adapters/mongodb";
import { db } from "./mongo.ts";
import { magicLink, jwt } from "better-auth/plugins";
import { sendEmail } from "./postmark.ts"

export const auth = betterAuth({
	database: mongodbAdapter(db),
	plugins: [
        magicLink({
            sendMagicLink: async ({ email, url }, request) => {
            	await sendEmail(email, url);
            },
            expiresIn: 300,
            disableSignUp: false,
        }),
        jwt(),
    ]
});
