/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    files: [
      "./views/**/*.{templ,html}", // All Templ and HTML files in views/
      "./static/cards.js",
      "./static/email.js",
      "./static/animation.js",
    ],
    exclude: [
      "**/node_modules/**",
      "**static/encom-globe.js",
      "**static/data.js",
    ],
  },
  theme: {
    extend: {},
  },
  plugins: [],
};
