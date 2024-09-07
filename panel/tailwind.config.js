/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/node_modules/preline/dist/*.js", "./views/**/*.templ"],
  theme: {
    container: {
      padding: "2rem",
      center: true,
      screens: {
        sm: "540px",
        md: "650px",
        lg: "900px",
        xl: "1100px",
        "2xl": "1300x",
      },
    },
  },
};
