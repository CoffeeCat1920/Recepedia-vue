export default {
  darkMode: "media",
  future: {
    // removeDeprecatedGapUtilities: true,
    // purgeLayersByDefault: true,
  },
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  purge: [],
  theme: {
    extend: {
      colors: {
        background_primary: "#F4E5CE", 
        background_secondary: "#E4D29D",
        foreground_primary: "#000000", 
        foreground_secondary: "#77653E"
      },    
      fontFamily: {
        monomakh: ['Monomakh', 'sans-serif'],
        roboto: ['Roboto', 'sans-serif'],
        robotoMono: ['Roboto Mono', 'monospace'],
      },
    },
  },
  variants: {},
  plugins: [],
}
