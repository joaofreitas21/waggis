document.addEventListener("DOMContentLoaded", function () {
    let currentKey = null;
  
    document.body.addEventListener("click", async function (e) {
      const tech = e.target.closest("[data-key]");
      if (tech) {
        const key = tech.getAttribute("data-key");
        const colorClass = tech.getAttribute("data-color") || "";
        if (!key) return;
  
        const slot = document.getElementById("card-slot");
  
        // Toggle close if same item and visible
        if (currentKey === key && !slot.classList.contains("hidden")) {
          closeCard();
          currentKey = null;
          return;
        }
  
        await showCardIntoSlot(key, colorClass);
        currentKey = key;
      }
    });
  });
  
  async function showCardIntoSlot(key, colorClass) {
    const slot = document.getElementById("card-slot");
    const rest = document.getElementById("info-line");
    const globe = document.getElementById("globe-wrap");
  
    // Apply color (text-* controls border-current/decoration-current inside the card)
    slot.classList.forEach(c => {
      if (c.startsWith("text-")) slot.classList.remove(c);
    });
    if (colorClass) slot.classList.add(colorClass);
  
    const resp = await fetch("/card?key=" + encodeURIComponent(key));
    slot.innerHTML = resp.ok
      ? await resp.text()
      : "<div class='text-red-400 p-6'>Failed to load card :(</div>";
  
    slot.classList.remove("hidden");
    if (rest) rest.classList.add("hidden");
    if (globe) globe.classList.add("hidden");
    document.body.style.overflow = "hidden";
  }
  
  function closeCard() {
    const slot = document.getElementById("card-slot");
    const rest = document.getElementById("info-line");
    const globe = document.getElementById("globe-wrap");
  
    slot.innerHTML = "";
    slot.classList.add("hidden");
    if (rest) rest.classList.remove("hidden");
    if (globe) globe.classList.remove("hidden");

    document.body.style.overflow ="";
  }