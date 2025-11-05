const anim = [
    { t: "{ }", ms: 200 },
    { t: "{_}", ms: 200 },
    { t: "{ }", ms: 200 },
    { t: "{_}", ms: 200 },
    { t: "{W_}", ms: 100 },
    { t: "{WA_}", ms: 100 },
    { t: "{WAG_}", ms: 100 },
    { t: "{WAGG_}", ms: 100 },
    { t: "{WAGGI_}", ms: 100 },
    { t: "{WAGGIS_}", ms: 100 },
    { t: "{WAGGIS }", ms: 200 },
    { t: "{WAGGIS_}", ms: 200 },
    { t: "{WAGGIS }", ms: 200 },
    { t: "{WAGGIS_}", ms: 200 },
    { t: "{WAGGIS}", ms: 200 },
    { t: "{WAGGIS}", ms: 200 }
];

let i = 0;
const text = document.getElementById('waggis-text');

function step() {
    text.textContent = anim[i].t;
    if (i < anim.length - 1) {
        i++;
        setTimeout(step, anim[i].ms);
    } else {
        // Animation done, trigger transition
        setTimeout(moveToTop, 500);
    }
}
if (text) step();

function moveToTop() {
    text.classList.remove('text-6xl');
    text.classList.add('text-3xl');
    text.style.transition = "all 0.7s cubic-bezier(0.4,0,0.2,1)";
    text.style.transform = "translateY(-45vh)";

    setTimeout(() => {
        const mainContent = document.getElementById('main-content');
        if (mainContent) {
            mainContent.classList.remove('opacity-0');
            mainContent.classList.add('opacity-100');
        }

        if(typeof showGlobe === 'function') {
            setTimeout(() => {
                showGlobe();
            }, 500);
        }
        
    }, 700);
}
