/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go",],
  safelist: [],
  plugins: [require("daisyui")],
  // daisyui: {
  //   themes: ["valentine"],
  // },
}