const anim = [
    { t: "{ }", ms: 200 },
    { t: "{_}", ms: 200 },
    { t: "{ }", ms: 200 },
    { t: "{_}", ms: 200 },
    { t: "{N_}", ms: 200 },
    { t: "{NO_}", ms: 200 },
    { t: "{NOT_}", ms: 200 },
    { t: "{NOT _}", ms: 200 },
    { t: "{NOT F_}", ms: 200 },
    { t: "{NOT FO_}", ms: 200 },
    { t: "{NOT FOU_ }", ms: 200 },
    { t: "{NOT FOUN_}", ms: 200 },
    { t: "{NOT FOUND_}", ms: 200 },
    { t: "{NOT FOUND }", ms: 200 },
    { t: "{NOT FOUND}", ms: 200 },
];

let i = 0;
const text = document.getElementById('notfound-text');

function step() {
    text.textContent = anim[i].t;
    if (i < anim.length - 1) {
        i++;
        setTimeout(step, anim[i].ms);
    }
}
if (text) step();


