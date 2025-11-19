import postmark from 'postmark';

var client = new postmark.ServerClient(import.meta.env.POSTMARK_SERVER_TOKEN);

export async function sendEmail(to: string, magicLink: string) {
	try {
		await client.sendEmail({
			From: import.meta.env.POSTMARK_FROM_ADDRESS,
			To: to,
			Subject: "Your magic link",
			TextBody: `<a href="${magicLink}">Magic</a>`
		})

	} catch (err) {
		console.error("Failed to send magic link:", err)
		throw err;
	}
}
