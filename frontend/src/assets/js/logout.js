const logoutBtn = document.getElementById('logout-btn');

logoutBtn.addEventListener('click', async () => {
	try {
		const res = await fetch("/api/auth/sign-out", {
			method: "POST",
		});

		if (res.ok) {
			window.location.href = "/";
		} else {
			console.error("Logout failed");
		}
	} catch (err) {
		console.error(err);
	}
})
