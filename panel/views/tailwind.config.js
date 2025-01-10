/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.templ"],
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
    extend: {
      fontFamily: {
        sans: ["Geist", "Inter", "system-ui", "sans-serif", "ui-sans-serif"],
        serif: ["Instrument Serif", "ui-serif", "serif"],
      },
      colors: {
        light: {
          0: "#f4f4f4",
          10: "#e9e9e9",
          20: "#dedede",
        },
        red: {
          0: "#db7070",
          10: "#c76262",
          20: "#b35454",
        },
      },
    },
  },
};
