/** @type {import('tailwindcss').Config} */
export default {
    content: ["./templates/*.html", "./static/*.js"],
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

