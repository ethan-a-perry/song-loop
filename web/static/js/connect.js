const connectBtn = document.getElementById("connect");

connectBtn.addEventListener("click", () => {
    window.location.href = "/api/spotify/connect";
});
