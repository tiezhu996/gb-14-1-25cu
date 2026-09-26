/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#eef6ff',
          100: '#d9eaff',
          200: '#bcdbff',
          300: '#8ec5ff',
          400: '#59a4ff',
          500: '#3380ff',
          600: '#1d60f5',
          700: '#164be1',
          800: '#193eb6',
          900: '#1a398f',
        },
      },
    },
  },
  plugins: [],
}
