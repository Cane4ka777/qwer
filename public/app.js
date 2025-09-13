(() => {
  const $ = (id) => document.getElementById(id);
  const statusEl = $("status");
  const metrics = $("metrics");
  const mcount = $("mcount");
  const acount = $("acount");
  const scount = $("scount");
  const versionEl = document.getElementById("api-version");
  const yearEl = document.getElementById("year");
  yearEl && (yearEl.textContent = new Date().getFullYear());

  async function pingAPI() {
    try {
      const res = await fetch("/api", { headers: { "Accept": "application/json" } });
      if (!res.ok) throw new Error("API not ok");
      const data = await res.json();
      statusEl.textContent = "API Online";
      statusEl.classList.add("ok");
      if (versionEl && data && data.data && data.data.version) {
        versionEl.textContent = data.data.version;
      }
    } catch (e) {
      statusEl.textContent = "API Offline";
      statusEl.classList.add("err");
      console.error(e);
    }
  }

  async function loadBand() {
    try {
      const res = await fetch("/api/band");
      if (!res.ok) throw new Error("band not ok");
      const json = await res.json();
      const band = json.data || {};
      const members = Array.isArray(band.members) ? band.members.length : 0;
      const albums = Array.isArray(band.discography) ? band.discography.length : 0;
      let songs = 0;
      if (Array.isArray(band.discography)) {
        for (const a of band.discography) {
          if (Array.isArray(a.songs)) songs += a.songs.length;
        }
      }
      mcount.textContent = members;
      acount.textContent = albums;
      scount.textContent = songs;
      metrics.hidden = false;
    } catch (e) {
      console.error(e);
    }
  }

  pingAPI();
  loadBand();
})();
