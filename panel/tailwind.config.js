/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ"],
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
        sans: ["Inter", "system-ui", "sans-serif", "ui-sans-serif"],
        serif: ["Instrument Serif", "ui-serif", "serif"],
      },
      colors: {
        light: {
          0: "#f4f4f4",
          10: "#e9e9e9",
          20: "#dedede",
        },
        red: "#db7070",
        green: "#7c9f4b",
        yellow: "#d69822",
        blue: {
          0: "#6587bf",
          20: "#516c99",
        },
        magenta: "#b870ce",
        cyan: "#509c93",
      },
    },
  },
};
