/**
 * @type {import('tailwindcss').Config}
 */
module.exports = {
  content: ["./internal/**/*.templ"],
  safelist: ["alert-error", "alert-success", "alert-warning", "alert-info"],
  theme: {
    extend: {
      animation: {
        "slide-in": "slide-in 0.5s ease-out forwards",
      },
      keyframes: {
        "slide-in": {
          "0%": { transform: "translateX(100%)", opacity: "0" },
          "100%": { transform: "translateX(0)", opacity: "1" },
        },
      },
    },
  },
};
