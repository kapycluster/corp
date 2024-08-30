/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/node_modules/preline/dist/*.js", "./views/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("./views/node_modules/preline/plugin.js")],
};
