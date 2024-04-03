/** @type {import('tailwindcss').Config} */
export default {
    content: ["./internal/**/*.templ", "./static/**/*.js"],
    theme: {
        extend: {
            colors: {
                primary: "#d8cb11"
            },
            fontFamily: {
                title: "'Madimi One', sans-serif",
                body: "Inter, sans-serif"
            }
        },
    },
    plugins: [],
}

