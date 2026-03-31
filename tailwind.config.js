/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    files: [
      "./views/**/*.{templ,html}", 
      "./static/cards.js",
      "./static/email.js",
      "./static/animation.js",
      "./**/*.go",
      "./static/ctf.js"
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
  plugins: [
    require('@tailwindcss/typography'),
  ],
};
