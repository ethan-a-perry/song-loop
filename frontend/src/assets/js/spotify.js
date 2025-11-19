import { authClient } from '../../lib/auth-client'

export async function connectSpotify() {
	try {
		const { data, error } = await authClient.token()
		if (error) {
	  	console.error("Could not access JWT")
				return
		}

		const res = await fetch("http://localhost:8080/api/spotify/connect", {
			method: "POST",
			headers: { "Authorization": `Bearer ${data.token}` }
		})

		if (res.ok) {
			console.log("connect spotify success")
			const data = await res.json()
      // Redirect user to Spotify
      window.location.href = data.authUrl
  	}
	  else {
	  	console.log("connect spotify fail")
	  }
	} catch (err) {
		console.error(err)
	}
}
