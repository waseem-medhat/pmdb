/** @type {import('tailwindcss').Config} */
export default {
    content: ["./internal/templs/**/*.templ", "./static/**/*.js"],
    theme: {
        extend: {
            colors: {
                primary: "#7cdd2c"
            },
            fontFamily: {
                title: "'Madimi One', sans-serif",
                body: "Inter, sans-serif"
            }
        },
    },
    plugins: [],
}

