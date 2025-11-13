const form = document.getElementById('song-loop-form');
const statusMessage = document.getElementById('form-status');

form.addEventListener('submit', async (e) => {
	console.log("tried")
	e.preventDefault();

	statusMessage.textContent = "Sending...";

	const formData = new FormData(form);

	try {
		const res = await fetch('http://127.0.0.1:8080/loop', {
			method: 'POST',
			body: JSON.stringify({
				start: parseInt(formData.get('start')),
				end: parseInt(formData.get('end')),
			}),
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (res.ok) {
			statusMessage.textContent = "Loop has been set";
		}
		else {
			statusMessage.textContent = "Loop has not been set";
		}
	} catch (err) {
		statusMessage.textContent = "Loop has not been set";
		console.error(err);
	}
})
