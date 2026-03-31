"use strict";

let activeSearch = "";
let activeTag    = "";
let debounceTimer = null;

const grid          = document.getElementById("grid");
const resultsCount  = document.getElementById("resultsCount");
const input         = document.getElementById("searchInput");
const activeTagPill = document.getElementById("activeTagPill");
const activeTagLabel = document.getElementById("activeTagLabel");
const clearTagBtn   = document.getElementById("clearTag");


function renderPrompt() {
  resultsCount.textContent = "";
  grid.innerHTML = `
    <div class="col-span-full text-center text-white/30 font-mono text-sm py-16">
      Start typing to search writeups<span class="animate-blink"></span>
    </div>
  `;
}

function updateTagPill() {
  if (activeTag) {
    activeTagLabel.textContent = activeTag;
    activeTagPill.classList.remove("hidden");
    activeTagPill.classList.add("inline-flex");
  } else {
    activeTagPill.classList.add("hidden");
    activeTagPill.classList.remove("inline-flex");
  }
}


function fetchResults() {
  const q   = activeSearch.trim();
  const tag = activeTag.trim();

  updateTagPill();

  if (!q && !tag) {
    renderPrompt();
    return;
  }

  const params = new URLSearchParams({ q, tag });

  fetch(`/api/ctf/search?${params}`)
    .then((res) => {
      if (!res.ok) throw new Error(res.statusText);
      return res.json();
    })
    .then((data) => {
      resultsCount.textContent = `> ${data.count} result${data.count !== 1 ? "s" : ""}`;
      grid.innerHTML = data.html || `
        <div class="col-span-full text-center text-white/30 font-mono text-sm py-16">
          // no writeups match your query
        </div>
      `;
    })
    .catch(() => {
      resultsCount.textContent = "";
      grid.innerHTML = `
        <div class="col-span-full text-center text-white/30 font-mono text-sm py-16">
          // failed to load results
        </div>
      `;
    });
}



input.addEventListener("input", (e) => {
  activeSearch = e.target.value;
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(fetchResults, 250);
});


grid.addEventListener("click", (e) => {
  const tagBtn = e.target.closest("[data-tag]");
  const cardEl = e.target.closest("[data-slug]");

  if (tagBtn) {
    e.stopPropagation();
    const clicked = tagBtn.getAttribute("data-tag") || "";
    // Toggle off if already active
    activeTag = activeTag === clicked ? "" : clicked;
    fetchResults();
    return;
  }

  if (cardEl) {
    const slug = cardEl.getAttribute("data-slug");
    if (slug) window.location.href = `/ctf/report?slug=${encodeURIComponent(slug)}`;
  }
});


clearTagBtn.addEventListener("click", () => {
  activeTag = "";
  fetchResults();
});

renderPrompt();
