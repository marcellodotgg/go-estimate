/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.html"],
  theme: {
    extend: {
      borderColor: {
        DEFAULT: "rgba(0, 0, 0, 0.5)",
      },
      colors: {
        DEFAULT: "#222",
      }
    },
    container: {
      center: true,
      screens: {
        sm: "100%",
        md: "100%",
        lg: "100%",
        xl: "968px",
        "2xl": "968px",
      },
      padding: {
        DEFAULT: "1rem",
        xl: "2rem",
      },
    },
  },
};
