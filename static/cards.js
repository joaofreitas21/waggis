document.addEventListener("DOMContentLoaded", () => {
  let currentKey = null;

  document.body.addEventListener("click", async (e) => {
    // Mini-card click (open item overlay)
    const itemBtn = e.target.closest(".mini-card[data-item-id]");
    if (itemBtn) {
      const itemId = itemBtn.getAttribute("data-item-id");
      if (itemId && currentKey) {
        await showItemOverlay(currentKey, itemId);
      }
      return;
    }

    // Close overlay button
    if (e.target.classList.contains("item-overlay-close")) {
      closeItemCardOverlay();
      return;
    }

    // Prevent overlay wrapper clicks from toggling the main card
    if (e.target.closest("#item-card-overlay")) {
      e.stopPropagation();
      return;
    }

    // Tech card click (open/close main card)
    const tech = e.target.closest("[data-key]");
    if (tech) {
      const key = tech.getAttribute("data-key");
      const colorClass = tech.getAttribute("data-color") || "";
      if (!key) return;

      const slot = document.getElementById("card-slot");

      // Toggle close if same card is already open
      if (currentKey === key && slot && !slot.classList.contains("hidden")) {
        closeCard();
        currentKey = null;
        return;
      }

      await showCardIntoSlot(key, colorClass);
      currentKey = key;
    }
  });

  window.closeItemCardOverlay = closeItemCardOverlay;
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

async function showItemOverlay(key, itemId) {
  const overlay = document.getElementById("item-card-overlay");
  if (!overlay) return;

  const resp = await fetch(
    `/card/item?key=${encodeURIComponent(key)}&id=${encodeURIComponent(itemId)}`
  );

  if (resp.ok) {
    overlay.innerHTML = await resp.text();
  } else {
    overlay.innerHTML = "<div class='text-red-400 p-6'>Failed to load card :(</div>";
  }

  overlay.classList.remove("hidden");
  overlay.addEventListener("click", (evt) => evt.stopPropagation(), { once: true });
  document.body.style.overflow = "hidden";
}

function closeItemCardOverlay() {
  const overlay = document.getElementById("item-card-overlay");
  if (overlay) {
    overlay.innerHTML = "";
    overlay.classList.add("hidden");
  }
  document.body.style.overflow = "";
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