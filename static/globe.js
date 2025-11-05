let globeInstance = null;

function initGlobe() {
  const container = document.getElementById("globe-container");
  if (!container) return;

  if (typeof ENCOM === "undefined" || typeof ENCOM.Globe ==="undefined" ) {
    console.error("Encom Globe library not found");
    return;
  }

  if (typeof grid === "undefined" || !grid.tiles){
    console.error("grid library not found");
    return;
  }

  const width = container.clientWidth || 800;
  const height = container.clientHeight || 600;

  globeInstance = new ENCOM.Globe(width, height, {
    font: "Inconsolata",
    data: [], 
    tiles: grid.tiles, 
    baseColor: "#000000",      // Black background
    markerColor: "#FFD700",    // Red markers 
    pinColor: "#FFD700",       // Golden pins
    satelliteColor: "#FF0000", // satellites
    scale: 1.2,
    dayLength: 14000,
    introLinesDuration: 2000,
    maxPins: 500,
    maxMarkers: 10,
    viewAngle: 0.1
  });

  container.appendChild(globeInstance.domElement);

  globeInstance.init(function() {
    console.log("Globe started");
    animate();

    setTimeout(() => {
        addRandomSatellites(4);
    }, 1000);
  });
}

function animate() {
    if(globeInstance){
        globeInstance.tick();
        requestAnimationFrame(animate);
    }
}

function showGlobe() {
  const container = document.getElementById("globe-container");
  if (container) {
    container.classList.remove("opacity-0");
    container.classList.add("opacity-100");

    setTimeout(() => {
      initGlobe();

      setTimeout(() => {
        fetchUserMarker();
      }, 500);
    }, 100);
  }
}

function addMarker(latitude, longitude, label) {
  if (!globeInstance) return;

  globeInstance.addMarker(latitude,longitude,label,false);
}

async function fetchUserMarker() {
    try {
        const response = await fetch("/api/ip");
        if(!response.ok) {
            throw new Error("Failed to fetch user marker IP data");
        }
        const data = await response.json();

        if(data.longitude && data.latitude) {
            const label = data.ip || "Unknown IP Address"
            addMarker(data.latitude, data.longitude, label);
            console.log(`User marker added at${data.latitude}, ${data.longitude} with IP: ${label}`);
        }
        
    } catch (error) {
        console.error("Error fetching user marker:", error);
    }
}

function addRandomSatellites(count){
    if(!globeInstance) return;

    //8 satellites position added - enable as needed
    const positions = [
        { lat: 45.0, lon: -75.0 },   // North America (upper-left quadrant)
        { lat: 50.0, lon: 10.0 },     // Europe (upper-right quadrant)
        { lat: -25.0, lon: 135.0 },   // Australia (lower-right quadrant)
        { lat: -20.0, lon: -45.0 },   // South America (lower-left quadrant)
        { lat: 30.0, lon: 80.0 },     // Asia (upper-right)
        { lat: -30.0, lon: 25.0 },    // Southern Africa (lower-center)
        { lat: 60.0, lon: -120.0 },   // Northern Canada (upper-left)
        { lat: 0.0, lon: -150.0 },    // Pacific (center-left)
    ];

    for (let i = 0; i < count; i++){

        let lat, lon;

        if (i < positions.length){

            const pos = positions[i];
            lat = pos.lat + (Math.random() * 20) - 10;
            lon = pos.lon + (Math.random() * 20) - 10;
        }

        const altitude = 1.0;

        const options = {
            coreColor: "#FF0000",
            numWaves: 8,
            shieldColor: "#FFFFFF"
        };

        globeInstance.addSatellite(lat,lon,altitude,options);
        console.log(`Satellite ${i + 1} added at ${lat.toFixed(2)}, ${lon.toFixed(2)} with altitude ${altitude.toFixed(2)}`);
    }
}

if (typeof module !== "undefined" && module.exports) {
  module.exports = { showGlobe, initGlobe, addMarker, addRandomSatellites};
}


