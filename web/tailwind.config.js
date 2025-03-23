export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--background)",
        text: "var(--text)",
        headerBackground: "var(--header-background)",
      },
      fontFamily: {
        roboto: ['Roboto', 'sans-serif'],
        monomakh: ['Monomakh', 'serif'],
      },
      typography: {
        DEFAULT: {
          css: {
            body: {
              fontFamily: 'Roboto, sans-serif',
              fontWeight: 'normal',
              backgroundColor: 'var(--background)',
              color: 'var(--text)',
              margin: '0',
              padding: '0',
            },
            h1: {
              fontFamily: 'Monomakh, serif',
              textAlign: 'center',
            },
            h2: {
              fontFamily: 'Monomakh, serif',
              fontWeight: 'normal',
              textAlign: 'center',
              padding: '20px',
            },
            h3: {
              fontFamily: 'Monomakh, serif',
              textAlign: 'center',
            },
            h4: {
              fontFamily: 'Monomakh, serif',
              textAlign: 'center',
            },
            h5: {
              fontFamily: 'Monomakh, serif',
              textAlign: 'center',
            },
            h6: {
              fontFamily: 'Monomakh, serif',
              textAlign: 'center',
            },
            ".header": {
              padding: "20px",
              backgroundColor: "var(--header-background)",
              borderRadius: "0 0 240px 100% / 240px",
            },
          },
        },
      },
    },
  },
  plugins: [require('@tailwindcss/typography')],
};
